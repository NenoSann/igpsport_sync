package igpsportsync

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

// loadEnvFile loads environment variables from .env file
func loadEnvFile(filepath string) (map[string]string, error) {
	envVars := make(map[string]string)
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			value = strings.Trim(value, "\"'")
			envVars[key] = value
		}
	}

	return envVars, scanner.Err()
}

// TestIntegration_BusinessFlow tests the complete business workflow:
// 1. Load environment variables from .env file
// 2. Login with credentials
// 3. GetActivityList to retrieve activities
// 4. GetActivityDownloadUrl for the first activity
// 5. DownloadFile to download activity data
//
// Each step is dependent on the previous step's success.
// If any step fails, the entire test fails and subsequent steps are skipped.
func TestIntegration_BusinessFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// ===== Step 1: Load environment variables from .env file =====
	t.Log("Step 1: Loading environment variables from .env file...")
	envVars, err := loadEnvFile("./.env")
	if err != nil {
		t.Fatalf("Step 1 FAILED: Could not load .env file: %v", err)
	}

	username, ok := envVars["username"]
	if !ok || username == "" {
		t.Fatalf("Step 1 FAILED: Username not found or empty in .env file")
	}

	password, ok := envVars["password"]
	if !ok || password == "" {
		t.Fatalf("Step 1 FAILED: Password not found or empty in .env file")
	}

	t.Logf("Step 1 PASSED: Successfully loaded credentials (username: %s)", username)

	// ===== Step 2: Login with credentials =====
	t.Log("Step 2: Attempting to login...")
	config := Config{
		Username: username,
		Password: password,
		PageSize: 20,
	}

	igpsport := &IgpsportSync{
		Config: config,
	}

	err = igpsport.init(config)
	if err != nil {
		t.Fatalf("Step 2 FAILED: Login failed: %v", err)
	}

	if igpsport.LoginResult == nil || igpsport.LoginResult.Access_token == "" {
		t.Fatalf("Step 2 FAILED: Failed to obtain access token after login")
	}

	accessToken := igpsport.LoginResult.Access_token
	t.Logf("Step 2 PASSED: Successfully logged in, access token: %s...", accessToken[:min(20, len(accessToken))])

	// ===== Step 3: GetActivityList to retrieve activities =====
	t.Log("Step 3: Retrieving activity list...")
	activityListResp, err := igpsport.GetActivityList(1, GPX)
	if err != nil {
		t.Fatalf("Step 3 FAILED: Could not retrieve activity list: %v", err)
	}

	if activityListResp == nil || activityListResp.Data.Rows == nil || len(activityListResp.Data.Rows) == 0 {
		t.Fatalf("Step 3 FAILED: No activities found in the response")
	}

	activities := activityListResp.Data.Rows
	t.Logf("Step 3 PASSED: Retrieved %d activities from page 1 (total pages: %d)", len(activities), activityListResp.Data.TotalPage)

	// Display activity details
	for i, activity := range activities {
		t.Logf("  Activity %d: ID=%d, Title=%s, StartTime=%s", i+1, activity.RideID, activity.Title, activity.StartTime)
	}

	// Use the first activity for subsequent steps
	firstActivity := activities[0]
	t.Logf("Using first activity for download test: ID=%d, Title=%s", firstActivity.RideID, firstActivity.Title)

	// ===== Step 4: GetActivityDownloadUrl for the first activity =====
	t.Logf("Step 4: Retrieving download URL for activity ID=%d...", firstActivity.RideID)
	downloadURL, err := igpsport.GetActivityDownloadUrl(firstActivity.RideID)
	if err != nil {
		t.Fatalf("Step 4 FAILED: Could not retrieve download URL: %v", err)
	}

	if downloadURL == nil || *downloadURL == "" {
		t.Fatalf("Step 4 FAILED: Download URL is empty or nil")
	}

	t.Logf("Step 4 PASSED: Successfully retrieved download URL: %s", *downloadURL)

	// ===== Step 5: DownloadFile to download activity data =====
	t.Logf("Step 5: Downloading activity file...")
	fileData, err := igpsport.DownloadFile(*downloadURL, fmt.Sprintf("%d", firstActivity.RideID), GPX)
	if err != nil {
		t.Fatalf("Step 5 FAILED: Could not download file: %v", err)
	}

	if fileData == nil || len(fileData) == 0 {
		t.Fatalf("Step 5 FAILED: Downloaded file is empty")
	}

	t.Logf("Step 5 PASSED: Successfully downloaded %d bytes of data", len(fileData))

	// ===== All steps completed successfully =====
	t.Log("")
	t.Log("============================================================")
	t.Log("✓ BUSINESS FLOW TEST COMPLETED SUCCESSFULLY")
	t.Log("============================================================")
	t.Logf("  1. ✓ Loaded credentials from .env")
	t.Logf("  2. ✓ Successfully authenticated (token: %s...)", accessToken[:min(20, len(accessToken))])
	t.Logf("  3. ✓ Retrieved activity list (%d activities)", len(activities))
	t.Logf("  4. ✓ Obtained download URL for activity: %s", firstActivity.Title)
	t.Logf("  5. ✓ Downloaded activity file (%d bytes)", len(fileData))
}

// TestIntegration_LoadEnvFile tests the loadEnvFile function
func TestIntegration_LoadEnvFile(t *testing.T) {
	envVars, err := loadEnvFile(".env")
	if err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	if username, ok := envVars["username"]; !ok || username == "" {
		t.Error("username not found or empty in .env file")
	} else {
		t.Logf("Loaded username: %s", username)
	}

	if password, ok := envVars["password"]; !ok || password == "" {
		t.Error("password not found or empty in .env file")
	} else {
		t.Logf("Loaded password: %s", strings.Repeat("*", len(password)))
	}
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

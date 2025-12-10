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

// TestIntegration_Login tests the login functionality with real server
func TestIntegration_Login(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Load environment variables
	envVars, err := loadEnvFile(".env")
	if err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	username, ok := envVars["username"]
	if !ok {
		t.Fatal("username not found in .env file")
	}

	password, ok := envVars["password"]
	if !ok {
		t.Fatal("password not found in .env file")
	}

	config := Config{
		Username: username,
		Password: password,
		PageSize: 20,
	}

	// Create instance with custom client to avoid auto-login
	igpsport := &IgpsportSync{
		Config: config,
	}

	// Manually initialize
	igpsport.init(config)

	// Verify access token was obtained
	if igpsport.AccessToken == "" {
		t.Error("Failed to obtain access token after login")
	} else {
		t.Logf("Successfully obtained access token: %s...", igpsport.AccessToken[:20])
	}
}

// TestIntegration_GetActivityList tests the GetActivityList functionality with real server
func TestIntegration_GetActivityList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Load environment variables
	envVars, err := loadEnvFile(".env")
	if err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	username, ok := envVars["username"]
	if !ok {
		t.Fatal("username not found in .env file")
	}

	password, ok := envVars["password"]
	if !ok {
		t.Fatal("password not found in .env file")
	}

	config := Config{
		Username: username,
		Password: password,
		PageSize: 20,
	}

	// Create instance and login
	igpsport := &IgpsportSync{
		Config: config,
	}
	igpsport.init(config)

	if igpsport.AccessToken == "" {
		t.Fatal("Failed to login before testing GetActivityList")
	}

	// Test GetActivityList with GPX format
	t.Run("GetActivityList_GPX", func(t *testing.T) {
		resp, err := igpsport.GetActivityList(1, GPX)
		if err != nil {
			t.Fatalf("GetActivityList failed: %v", err)
		}

		if resp == nil {
			t.Error("GetActivityList returned nil response")
		} else {
			t.Logf("Successfully retrieved activity list with %d rows", len(resp.Data.Rows))
			t.Logf("Total pages: %d", resp.Data.TotalPage)

			// Print activity details
			for i, row := range resp.Data.Rows {
				t.Logf("Activity %d: ID=%d, Name=%s", i+1, row.RideID, row.Title)
			}
		}
	})

	// Test GetActivityList with FIT format
	t.Run("GetActivityList_FIT", func(t *testing.T) {
		resp, err := igpsport.GetActivityList(1, FIT)
		if err != nil {
			t.Fatalf("GetActivityList failed: %v", err)
		}

		fmt.Printf("Response: %+v\n", resp)

		if resp == nil {
			t.Error("GetActivityList returned nil response")
		} else {
			t.Logf("Successfully retrieved activity list with %d rows (FIT format)", len(resp.Data.Rows))
		}
	})

	// Test GetActivityList with invalid page number
	t.Run("GetActivityList_InvalidPageNo", func(t *testing.T) {
		_, err := igpsport.GetActivityList(0, GPX)
		if err == nil {
			t.Error("Expected error for invalid page number, but got none")
		} else {
			t.Logf("Got expected error for invalid page number: %v", err)
		}
	})
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

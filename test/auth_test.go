package test

import (
	"strings"
	"testing"
)

// TestLoadEnvFile tests the LoadEnvFile function
func TestLoadEnvFile(t *testing.T) {
	envVars, err := LoadEnvFile("../.env")
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

// TestLogin tests the login functionality
func TestLogin(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := CreateTestClient(t)

	if client.LoginResult == nil || client.LoginResult.Access_token == "" {
		t.Fatalf("Failed to obtain access token after login")
	}

	accessToken := client.LoginResult.Access_token
	t.Logf("Successfully logged in, access token: %s...", accessToken[:Min(20, len(accessToken))])
}

// TestUserInfo tests the GetUserInfo functionality
func TestUserInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := CreateTestClient(t)

	userInfo, err := client.GetUserInfo()
	if err != nil {
		t.Fatalf("Failed to get user info: %v", err)
	}

	if userInfo == nil || userInfo.Data.NickName == "" {
		t.Fatalf("User info is empty or invalid")
	}

	t.Logf("User info retrieved successfully:")
	t.Logf("  NickName: %s", userInfo.Data.NickName)
	t.Logf("  MemberId: %d", userInfo.Data.MemberId)
	t.Logf("  RideNum: %d", userInfo.Data.RideNum)
	t.Logf("  RideDistance: %d", userInfo.Data.RideDistance)
}

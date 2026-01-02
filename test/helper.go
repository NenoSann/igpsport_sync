package test

import (
	"bufio"
	"os"
	"strings"
	"testing"

	igpsportsync "github.com/NenoSann/igpsport_sync"
)

// LoadEnvFile loads environment variables from .env file
func LoadEnvFile(filepath string) (map[string]string, error) {
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

// CreateTestClient creates a new IgpsportSync client for testing
func CreateTestClient(t *testing.T) *igpsportsync.IgpsportSync {
	t.Helper()

	envVars, err := LoadEnvFile("../.env")
	if err != nil {
		t.Fatalf("Could not load .env file: %v", err)
	}

	username, ok := envVars["username"]
	if !ok || username == "" {
		t.Fatalf("Username not found or empty in .env file")
	}

	password, ok := envVars["password"]
	if !ok || password == "" {
		t.Fatalf("Password not found or empty in .env file")
	}

	config := igpsportsync.Config{
		Username: username,
		Password: password,
	}

	client, err := igpsportsync.New(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	return client
}

// Min returns the smaller of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
package github

import (
	"os"
	"testing"
)

func TestBuildJWTToken(t *testing.T) {
	// Create temp key file
	tmpKeyFile := "/tmp/key.pem"
	os.WriteFile(tmpKeyFile, testPrivateKey, 0644)

	// Test Case 1: Invalid private key
	cfg := &GitHubAppConfig{
		ApplicationID: 123,
		PrivateKey:    "invalid",
	}
	ghApp := &GitHubApp{
		Config: cfg,
	}
	err := ghApp.buildJWTToken()
	if err == nil {
		t.Errorf("Expected error for invalid private key")
	}

	// Test Case 2: Valid private key
	cfg = &GitHubAppConfig{
		ApplicationID: 298674,
		PrivateKey:    tmpKeyFile,
	}
	ghApp = &GitHubApp{
		Config: cfg,
	}
	err = ghApp.buildJWTToken()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if ghApp.Auth.JWTToken == "" {
		t.Errorf("Empty JWT token")
	}

	// Cleanup
	os.Remove(tmpKeyFile)
}

package github

import (
	"testing"
)

func TestBuildJWTToken(t *testing.T) {
	// Test Case 1: Invalid private key
	cfg := &GitHubAppConfig{
		ApplicationID: "123",
		PrivateKey:    []byte("invalid"),
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
		ApplicationID: "298674",
		PrivateKey:    testPrivateKey,
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
}

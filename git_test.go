package github

import "testing"

func TestClone(t *testing.T) {
	// Test Case 1: Invalid repository URL
	cfg := &GitHubAppConfig{
		LocalPath: "/tmp/test_repo",
		RepoURL:   "invalid",
	}
	ghApp := &GitHubApp{
		Config: cfg,
	}
	err := ghApp.Clone()
	if err == nil {
		t.Errorf("Expected error for invalid repository URL")
	}
}

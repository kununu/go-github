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

	// // Test Case 2: Valid repository URL
	// cfg = &GitHubAppConfig{
	// 	LocalPath: "/tmp/test_repo",
	// 	RepoURL:   "https://github.com/book/git-test-repository",
	// }
	// ghApp = &GitHubApp{
	// 	Config: cfg,
	// }
	// err = ghApp.Clone()
	// if err != nil {
	// 	t.Errorf("Unexpected error: %v", err)
	// }
	// if ghApp.gitRepo == nil || ghApp.worktree == nil {
	// 	t.Errorf("Empty git repository or worktree")
	// }
}

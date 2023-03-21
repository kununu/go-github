package github

import (
	"errors"

	"github.com/go-git/go-git/v5"
)

// GitHubApp is a struct to interact with GitHub App's API
type GitHubApp struct {
	Config   *GitHubAppConfig
	Auth     Auth
	gitRepo  *git.Repository
	worktree *git.Worktree
}

// Configuration for the GitHub App interaction
type GitHubAppConfig struct {
	RepoURL        string
	ApplicationID  string
	InstallationID string
	LocalPath      string
	PrivateKey     []byte
}

var githubAPIURL string = "https://api.github.com"

// Creates a new GitHubApp struct configured by GitHubAppConfig
func NewGitHubApp(cfg *GitHubAppConfig) (*GitHubApp, error) {
	if cfg.ApplicationID == "" && cfg.InstallationID == "" {
		return nil, errors.New("provide App ID and Installation ID")
	}
	if cfg.LocalPath == "" {
		cfg.LocalPath = "./"
	}
	ghApp := &GitHubApp{
		Config: cfg,
	}
	err := ghApp.buildJWTToken()
	if err != nil {
		return nil, err
	}

	return ghApp, nil
}

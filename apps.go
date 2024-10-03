package github

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v54/github"
)

// GitHubApp is a struct to interact with GitHub App's API
type GitHubApp struct {
	Config       *GitHubAppConfig
	Auth         Auth
	gitClient    *git.Repository
	githubClient *github.Client
	worktree     *git.Worktree
}

// Configuration for the GitHub App interaction
type GitHubAppConfig struct {
	repoName       string
	RepoURL        string
	ApplicationID  int64
	InstallationID int64
	LocalPath      string
	PrivateKey     []byte
}

var githubAPIURL string = "https://api.github.com"

// extractRepoName extracts the repository name from the URL
func extractRepoName(gitURL string) (string, error) {
	parsedURL, err := url.Parse(gitURL)
	if err != nil {
		return "", err
	}

	// Splitting the path based on slashes
	parts := strings.Split(parsedURL.Path, "/")

	// Ensure at least one segment exists
	if len(parts) == 0 {
		return "", fmt.Errorf("could not extract repository name")
	}

	// The last part of the URL is typically 'repo.git', so we need to further process it to extract 'repo'
	repoName := parts[len(parts)-1]

	// Removing .git extension if present
	if strings.HasSuffix(repoName, ".git") {
		repoName = repoName[:len(repoName)-len(".git")]
	}

	return repoName, nil
}

// Creates a new GitHubApp struct configured by GitHubAppConfig
func NewGitHubApp(cfg *GitHubAppConfig) (*GitHubApp, error) {
	// check if application ID and installation ID are defined
	if cfg.ApplicationID == 0 && cfg.InstallationID == 0 {
		return nil, errors.New("provide App ID and Installation ID")
	}
	// set localPath if not passed
	if cfg.LocalPath == "" {
		cfg.LocalPath = "./"
	}

	var itr *ghinstallation.Transport
	var err error
	// get new key for github
	itr, err = ghinstallation.New(http.DefaultTransport, cfg.ApplicationID, cfg.InstallationID, cfg.PrivateKey)

	if err != nil {
		return nil, err
	}
	repoName, err := extractRepoName(cfg.RepoURL)
	if err != nil {
		return nil, err
	}
	// set repoName
	cfg.repoName = repoName

	// initialise a new GitHubApp
	ghApp := &GitHubApp{
		Config:       cfg,
		githubClient: github.NewClient(&http.Client{Transport: itr}),
	}

	// authenticate
	err = ghApp.authenticate()
	if err != nil {
		return nil, err
	}

	return ghApp, nil
}

func (ghApp *GitHubApp) authenticate() error {

	err := ghApp.buildJWTToken()
	if err != nil {
		return err
	}

	tkn, err := ghApp.GetAccessToken()
	if err != nil {
		return err
	}
	ghApp.Auth.Token = tkn

	return nil
}

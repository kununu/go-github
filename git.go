package github

import (
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// User information
type UserInfo struct {
	Name  string
	Email string
}

// Clones the repository
func (ghApp *GitHubApp) Clone() error {
	var err error
	ghApp.gitRepo, err = git.PlainClone(ghApp.Config.LocalPath, false, &git.CloneOptions{
		URL: ghApp.Config.RepoURL,
		Auth: &http.BasicAuth{
			Username: "github",
			Password: ghApp.Auth.Token,
		},
	})
	if err != nil {
		return err
	}
	ghApp.worktree, err = ghApp.gitRepo.Worktree()
	if err != nil {
		return err
	}

	return nil
}

// Check if repo has local changes
func (ghApp *GitHubApp) HasChanges() bool {
	status, err := ghApp.worktree.Status()
	if err != nil {
		return false
	}
	return status.IsClean()
}

// Adds files to the repository
func (ghApp *GitHubApp) Add(path string) {
	ghApp.worktree.Add(path)
}

// Commits the changes
func (ghApp *GitHubApp) Commit(msg string, user UserInfo) error {
	_, err := ghApp.worktree.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  user.Name,
			Email: user.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// Pushes the changes
func (ghApp *GitHubApp) Push() error {
	err := ghApp.gitRepo.Push(&git.PushOptions{
		RemoteName: "origin",
	})
	if err != nil {
		return err
	}

	return nil
}

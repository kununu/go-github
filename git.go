package github

import (
	"context"
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v70/github"
)

// User information
type UserInfo struct {
	Name  string
	Email string
}

// PullRequest information
type PullRequest struct {
	SourceBranch string
	TargetBranch string
}

// Clones the repository
func (ghApp *GitHubApp) Clone() error {
	var err error
	ghApp.gitClient, err = git.PlainClone(ghApp.Config.LocalPath, false, &git.CloneOptions{
		URL: ghApp.Config.RepoURL,
		Auth: &http.BasicAuth{
			Username: "github",
			Password: ghApp.Auth.Token,
		},
	})
	if err != nil {
		return err
	}
	ghApp.worktree, err = ghApp.gitClient.Worktree()
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
	return !status.IsClean()
}

// Adds files to the repository
func (ghApp *GitHubApp) Add(path string) error {
	_, err := ghApp.worktree.Add(path)
	return err
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
	err := ghApp.gitClient.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: "github",
			Password: ghApp.Auth.Token,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// Create a new branch
func (ghApp *GitHubApp) NewBranch(name string, checkout bool) error {

	return ghApp.worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", name)),
		Create: checkout,
	})

}

// Polling configuration for waiting on a PR's CI checks. Declared as vars
// (rather than consts) so tests can shrink them to keep polling loops fast.
var (
	checksPollInterval = 10 * time.Second
	checksWaitTimeout  = 10 * time.Minute
)

// Create new pull request
func (ghApp *GitHubApp) NewPullRequest(source, target, title, body string, merge bool) error {
	ctx := context.Background()

	// Create PR
	newPR := &github.NewPullRequest{
		Title:               github.Ptr(title),
		Head:                github.Ptr(source), // source branch
		Base:                github.Ptr(target), // target branch
		Body:                github.Ptr(body),
		MaintainerCanModify: github.Ptr(true),
	}

	pr, _, err := ghApp.githubClient.PullRequests.Create(ctx, "kununu", ghApp.Config.repoName, newPR)
	if err != nil {
		return err
	}

	if merge {

		if ghApp.Config.WaitForChecksToPass {
			err = waitForChecksToPass(ctx, ghApp.githubClient, "kununu", ghApp.Config.repoName, pr.GetHead().GetSHA())
			if err != nil {
				return err
			}
		}

		// Merge PR
		return mergePullRequest(ctx, ghApp.githubClient, "kununu", ghApp.Config.repoName, pr.GetNumber())
	}

	return nil

}

func waitForChecksToPass(ctx context.Context, client *github.Client, owner, repo, ref string) error {
	ctx, cancel := context.WithTimeout(ctx, checksWaitTimeout)
	defer cancel()

	ticker := time.NewTicker(checksPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout while waiting for checks to pass: %w", ctx.Err())
		case <-ticker.C:
		}

		result, _, err := client.Checks.ListCheckRunsForRef(ctx, owner, repo, ref, nil)
		if err != nil {
			return fmt.Errorf("listing check runs for %s: %w", ref, err)
		}

		pending := false
		for _, run := range result.CheckRuns {
			if run.GetStatus() != "completed" {
				pending = true
				continue
			}
			switch run.GetConclusion() {
			case "success", "skipped", "neutral":
			default:
				return fmt.Errorf("check %q did not pass: %s", run.GetName(), run.GetConclusion())
			}
		}

		if !pending && result.GetTotal() > 0 {
			return nil
		}
	}
}

func mergePullRequest(ctx context.Context, client *github.Client, owner, repo string, number int) error {
	opts := &github.PullRequestOptions{
		MergeMethod: "squash", // or "merge" or "rebase"
	}
	_, _, err := client.PullRequests.Merge(ctx, owner, repo, number, "Automated merge by bot", opts)
	return err
}

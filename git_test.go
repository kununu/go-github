package github

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/go-github/v70/github"
)

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

// newTestClient returns a *github.Client whose requests are directed at the
// given handler instead of the real GitHub API.
func newTestClient(t *testing.T, handler http.Handler) *github.Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client := github.NewClient(nil)
	baseURL, err := url.Parse(server.URL + "/")
	if err != nil {
		t.Fatalf("failed to parse test server URL: %v", err)
	}
	client.BaseURL = baseURL
	return client
}

// fastChecksPolling shrinks the polling interval/timeout for the duration of a
// test so wait loops complete in milliseconds, restoring the originals after.
func fastChecksPolling(t *testing.T, interval, timeout time.Duration) {
	t.Helper()
	origInterval, origTimeout := checksPollInterval, checksWaitTimeout
	checksPollInterval, checksWaitTimeout = interval, timeout
	t.Cleanup(func() {
		checksPollInterval, checksWaitTimeout = origInterval, origTimeout
	})
}

func TestWaitForChecksToPass(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		status      int
		timeout     time.Duration
		wantErr     bool
		wantTimeout bool // error must wrap context.DeadlineExceeded
		errContains string
	}{
		{
			name:   "all checks succeed",
			body:   `{"total_count":2,"check_runs":[{"name":"build","status":"completed","conclusion":"success"},{"name":"lint","status":"completed","conclusion":"success"}]}`,
			status: http.StatusOK,
		},
		{
			name:   "skipped and neutral count as passed",
			body:   `{"total_count":2,"check_runs":[{"name":"opt","status":"completed","conclusion":"skipped"},{"name":"info","status":"completed","conclusion":"neutral"}]}`,
			status: http.StatusOK,
		},
		{
			name:        "a failed check returns an error",
			body:        `{"total_count":1,"check_runs":[{"name":"build","status":"completed","conclusion":"failure"}]}`,
			status:      http.StatusOK,
			wantErr:     true,
			errContains: `check "build" did not pass`,
		},
		{
			// A timeout can surface either from the select branch or mid-flight
			// in the HTTP call, depending on where the deadline lands; both wrap
			// context.DeadlineExceeded.
			name:        "still-pending checks time out",
			body:        `{"total_count":1,"check_runs":[{"name":"build","status":"in_progress"}]}`,
			status:      http.StatusOK,
			timeout:     30 * time.Millisecond,
			wantErr:     true,
			wantTimeout: true,
		},
		{
			name:        "no checks registered times out",
			body:        `{"total_count":0,"check_runs":[]}`,
			status:      http.StatusOK,
			timeout:     30 * time.Millisecond,
			wantErr:     true,
			wantTimeout: true,
		},
		{
			name:        "API error is propagated",
			body:        `{"message":"boom"}`,
			status:      http.StatusInternalServerError,
			wantErr:     true,
			errContains: "listing check runs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeout := tt.timeout
			if timeout == 0 {
				timeout = time.Second
			}
			fastChecksPolling(t, time.Millisecond, timeout)

			handler := http.NewServeMux()
			handler.HandleFunc("GET /repos/kununu/test-repo/commits/{ref}/check-runs", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				w.Write([]byte(tt.body))
			})

			client := newTestClient(t, handler)
			err := waitForChecksToPass(context.Background(), client, "kununu", "test-repo", "deadbeef")

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.wantTimeout && !errors.Is(err, context.DeadlineExceeded) {
					t.Errorf("error %q does not wrap context.DeadlineExceeded", err.Error())
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestNewPullRequestSkipsChecksWhenDisabled(t *testing.T) {
	var checksCalled, mergeCalled bool

	handler := http.NewServeMux()
	handler.HandleFunc("POST /repos/kununu/test-repo/pulls", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"number":7,"head":{"sha":"abc123"}}`))
	})
	handler.HandleFunc("GET /repos/kununu/test-repo/commits/{ref}/check-runs", func(w http.ResponseWriter, r *http.Request) {
		checksCalled = true
		w.Write([]byte(`{"total_count":1,"check_runs":[{"name":"build","status":"completed","conclusion":"success"}]}`))
	})
	handler.HandleFunc("PUT /repos/kununu/test-repo/pulls/7/merge", func(w http.ResponseWriter, r *http.Request) {
		mergeCalled = true
		w.Write([]byte(`{"merged":true}`))
	})

	ghApp := &GitHubApp{
		Config:       &GitHubAppConfig{repoName: "test-repo", WaitForChecksToPass: false},
		githubClient: newTestClient(t, handler),
	}

	if err := ghApp.NewPullRequest("feature", "main", "title", "body", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if checksCalled {
		t.Errorf("check runs were polled even though WaitForChecksToPass is false")
	}
	if !mergeCalled {
		t.Errorf("PR was not merged")
	}
}

func TestNewPullRequestWaitsForChecksWhenEnabled(t *testing.T) {
	var checksCalled, mergeCalled bool

	handler := http.NewServeMux()
	handler.HandleFunc("POST /repos/kununu/test-repo/pulls", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"number":7,"head":{"sha":"abc123"}}`))
	})
	handler.HandleFunc("GET /repos/kununu/test-repo/commits/{ref}/check-runs", func(w http.ResponseWriter, r *http.Request) {
		checksCalled = true
		w.Write([]byte(`{"total_count":1,"check_runs":[{"name":"build","status":"completed","conclusion":"success"}]}`))
	})
	handler.HandleFunc("PUT /repos/kununu/test-repo/pulls/7/merge", func(w http.ResponseWriter, r *http.Request) {
		mergeCalled = true
		w.Write([]byte(`{"merged":true}`))
	})

	fastChecksPolling(t, time.Millisecond, time.Second)

	ghApp := &GitHubApp{
		Config:       &GitHubAppConfig{repoName: "test-repo", WaitForChecksToPass: true},
		githubClient: newTestClient(t, handler),
	}

	if err := ghApp.NewPullRequest("feature", "main", "title", "body", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !checksCalled {
		t.Errorf("check runs were not polled even though WaitForChecksToPass is true")
	}
	if !mergeCalled {
		t.Errorf("PR was not merged")
	}
}

func TestNewPullRequestSkipsMergeWhenNotRequested(t *testing.T) {
	var checksCalled, mergeCalled bool

	handler := http.NewServeMux()
	handler.HandleFunc("POST /repos/kununu/test-repo/pulls", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"number":7,"head":{"sha":"abc123"}}`))
	})
	handler.HandleFunc("GET /repos/kununu/test-repo/commits/{ref}/check-runs", func(w http.ResponseWriter, r *http.Request) {
		checksCalled = true
		w.Write([]byte(`{"total_count":1,"check_runs":[{"name":"build","status":"completed","conclusion":"success"}]}`))
	})
	handler.HandleFunc("PUT /repos/kununu/test-repo/pulls/7/merge", func(w http.ResponseWriter, r *http.Request) {
		mergeCalled = true
		w.Write([]byte(`{"merged":true}`))
	})

	ghApp := &GitHubApp{
		Config:       &GitHubAppConfig{repoName: "test-repo", WaitForChecksToPass: true},
		githubClient: newTestClient(t, handler),
	}

	if err := ghApp.NewPullRequest("feature", "main", "title", "body", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if checksCalled {
		t.Errorf("check runs were polled even though merge was not requested")
	}
	if mergeCalled {
		t.Errorf("PR was merged even though merge was not requested")
	}
}

# go-github

With the migration of GitHub to New Work's Single Sign-On (SSO), the GitHub user that was
being used in our tools (like Jenkins) was not beeing added as a new LDAP user.
With this in mind, we had the need to have a way of authenticating using
[GitHub Apps](https://docs.github.com/en/apps/creating-github-apps/creating-github-apps/about-apps).

The creation of this Go library is to be able to use in any code base that needs to authenticate
on GitHub, like our [fetch-bots-ips](https://github.com/kununu/fetch-bots-ips).

## Usage

### Get a new Access Token

To generate a new GitHub App access token that can be used to authenticate.

```go
package main

import "github.com/kununu/go-github/v1"

func main() {
// Create a new GithubApp with JWT authentication
ghApp, err := github.NewGitHubApp(&github.GitHubAppConfig{
		ApplicationID:  appId,
		InstallationID: instId,
		PrivateKey:     keyBytes,
})

token, err := ghApp.GetAccessToken()
if err != nil {
		panic(err)
}
```

## Installation

There is also a binary available that just outputs a GitHub token and can be
used in any tool by just setting the `GIT_ASKPASS` env variable.

```bash
GIT_ASKPASS=<path_to_binary>
```
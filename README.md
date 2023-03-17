# go-github

With the migration of GitHub to New Work's Single Sign-On (SSO), the GitHub user that was
being used in our tools (like Jenkins) was not beeing added as a new LDAP user.
With this in mind, we had the need to have a way of authenticating using
[GitHub Apps](https://docs.github.com/en/apps/creating-github-apps/creating-github-apps/about-apps).

The creation of this Go library is to be able to use in any code base that needs to authenticate
on GitHub, like our [fetch-bots-ips](https://github.com/kununu/fetch-bots-ips).

## Usage

To generate a new GitHub token that can be used to authenticate.

```go
package main

import "github.com/kununu/go-github/apps"

func main() {
  	// Create a new GithubApp with JWT authentication
	ctx, err := github.GetJWTContext(appId, keyBytes)
	if err != nil {
		panic(err)
	}

	// Get GitHub auth token for the specified installation
	token, err := github.GetAccessToken(ctx, instId)
	if err != nil {
		panic(err)
	}

}
```

## Installation

There is also a binary available that just outputs a GitHub token and can be
used in any tool by just setting the `GIT_ASKPASS` env variable.

```bash
GIT_ASKPASS=<path_to_binary>
```
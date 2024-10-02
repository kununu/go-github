# go-github

With the migration of our GitHub organisation to Single Sign-On (SSO), the GitHub user that was
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

import "github.com/kununu/go-github"

func main() {
	// Create a new GithubApp with JWT authentication
	ghApp, err := github.NewGitHubApp(&github.GitHubAppConfig{
		// GitHub APP ID to use to authenticate
		ApplicationID:  appId,
		// GitHub Installation ID for the APP
		InstallationID: instId,
		// GitHub APP generated private key path to file
		PrivateKeyFile: string,
		// GitHub APP generated private key value
		PrivateKey:     keyBytes,
	})

	token, err := ghApp.GetAccessToken()
	if err != nil {
		panic(err)
	}
}
```

## Provided binary

We also provide a binary that outputs a GitHub token and can be in any a linux command line by
setting the `GIT_ASKPASS` env variable to the binary.

The binary can be configured to be used byt setting the environment variables 
`GITHUB_APP_ID` `GITHUB_INST_ID` and `GITHUB_KEY_PATH` or `GITHUB_KEY_VALUE` or by passing the values using flags. 
Use the `-h` flag for more information on the available flags.

**NOTE:** If both GITHUB_KEY_PATH and GITHUB_KEY_VALUE are passed, only the first one is used.

```bash
GIT_ASKPASS=<path_to_binary>
```

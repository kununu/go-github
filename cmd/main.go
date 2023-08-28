package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/kununu/go-github"
)

var (
	appId  int64
	instId int64
	key    string
)

func init() {
	flag.StringVar(&key, "k", "", "Path to key file for authentication")
	flag.Int64Var(&appId, "a", 0, "App ID to use for authentication")
	flag.Int64Var(&instId, "i", 0, "Installation ID that identifies the APP installation ID on GitHub")
	flag.Parse()
}

func main() {
	// Get the values from the environment variables if they are set with parameters
	if appId == 0 {
		appIdInt, _ := strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
		appId = appIdInt
	}
	if instId == 0 {
		instIdInt, _ := strconv.ParseInt(os.Getenv("GITHUB_INST_ID"), 10, 64)
		instId = instIdInt
	}
	if key == "" {
		key = os.Getenv("GITHUB_KEY_PATH")
		if key == "" {
			// Read the key from STDIN
			stat, err := os.Stdin.Stat()
			if err != nil {
				fmt.Printf("error in stdin: %s", err)
				os.Exit(1)
			}
			if (stat.Mode() & os.ModeNamedPipe) == 0 {
				fmt.Printf("you need to pass the private key either with `-k` parameter or by setting GITHUB_KEY_PATH or even passing through STDIN\n")
				os.Exit(1)
			}
			key = "stdin"
		}
	}

	// Verify if the necessary information is set
	if appId == 0 || instId == 0 {
		fmt.Println("You need to define the App ID and the path to the key file")
		fmt.Println("by passing the values the -a, -i and -k options or")
		fmt.Println("by setting GITHUB_APP_ID, GITHUB_INST_ID and GITHUB_KEY_PATH environment variables.")
		os.Exit(1)
	}

	// Create a new GitHubApp
	ghApp, err := github.NewGitHubApp(&github.GitHubAppConfig{
		ApplicationID:  appId,
		InstallationID: instId,
		PrivateKey:     key,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Get GitHub auth token for the specified installation
	token, err := ghApp.GetAccessToken()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Printout GitHub Token
	fmt.Fprintf(os.Stdout, token)

}

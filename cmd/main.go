package main

import (
	"fmt"
	"os"

	"github.com/kununu/go-github"
)

func main() {
	// Parse the entry parameters
	config, err := github.ParseParameters()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create a new GitHubApp
	ghApp, err := github.NewGitHubApp(config)
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
	fmt.Fprintf(os.Stdout, "%s", token)

}

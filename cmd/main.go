package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/kununu/go-github/apps"
)

var (
	appId  int
	instId int
	key    string
)

func init() {
	flag.IntVar(&appId, "a", 0, "App ID to use for authentication")
	flag.StringVar(&key, "k", "", "Path to key file for authentication")
	flag.IntVar(&instId, "i", 0, "Installation ID that identifies the APP installation ID on GitHub")
	flag.Parse()
}

func main() {
	// Get the values from the environment variables if they are set with parameters
	if appId == 0 {
		appId, _ = strconv.Atoi(os.Getenv("GITHUB_APP_ID"))
	}
	if instId == 0 {
		instId, _ = strconv.Atoi(os.Getenv("GITHUB_INSTALLATION_ID"))
	}
	if key == "" {
		key = os.Getenv("GITHUB_KEY_PATH")
	}

	// Verify if the necessary information is set
	if appId == 0 || key == "" || instId == 0 {
		fmt.Println("You need to define the App ID and the path to the key file")
		fmt.Println("by passing the values the -a, -i and -k options or")
		fmt.Println("by setting GITHUB_APP_ID, GITHUB_INSTALLATION_ID and GITHUB_KEY_PATH environment variables.")
		os.Exit(0)
	}

	// Read the key from the file
	keyBytes, err := os.ReadFile(key)
	if err != nil {
		fmt.Println("error reading the key file")
		os.Exit(0)
	}

	// Create a new GithubApp with JWT authentication
	ctx, err := apps.GetJWTContext(appId, keyBytes)
	if err != nil {
		panic(err)
	}

	// Get GitHub auth token for the specified installation
	token, err := apps.GetAccessToken(ctx, instId)
	if err != nil {
		panic(err)
	}

	// Printout GitHub Token
	fmt.Println(token)

}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kununu/go-github"
)

var (
	appId  string
	instId string
	key    string
)

func init() {
	flag.StringVar(&appId, "a", "", "App ID to use for authentication")
	flag.StringVar(&key, "k", "", "Path to key file for authentication")
	flag.StringVar(&instId, "i", "", "Installation ID that identifies the APP installation ID on GitHub")
	flag.Parse()
}

func main() {
	keyBytes := []byte{}
	// Get the values from the environment variables if they are set with parameters
	if appId == "" {
		appId = os.Getenv("GITHUB_APP_ID")
	}
	if instId == "" {
		instId = os.Getenv("GITHUB_INST_ID")
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
			var lines []string
			reader := bufio.NewReader(os.Stdin)
			for {
				// read line from stdin using newline as separator
				line, _ := reader.ReadString('\n')
				// if line is empty, break the loop
				if len(strings.TrimSpace(line)) == 0 {
					break
				}
				//append the line to a slice
				lines = append(lines, line)
				keyBytes = append(keyBytes, []byte(line)...)
			}
			key = "stdin"
		}
	}

	// Verify if the necessary information is set
	if appId == "" || instId == "" {
		fmt.Println("You need to define the App ID and the path to the key file")
		fmt.Println("by passing the values the -a, -i and -k options or")
		fmt.Println("by setting GITHUB_APP_ID, GITHUB_INST_ID and GITHUB_KEY_PATH environment variables.")
		os.Exit(0)
	}

	// Read the key from the file
	if key != "stdin" {
		var err error
		keyBytes, err = os.ReadFile(key)
		if err != nil {
			fmt.Println("error reading the key file")
			os.Exit(1)
		}
	}

	// Create a new GitHubApp
	ghApp, err := github.NewGitHubApp(&github.GitHubAppConfig{
		ApplicationID:  appId,
		InstallationID: instId,
		PrivateKey:     keyBytes,
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

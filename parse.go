package github

import (
	"errors"
	"flag"
	"io"
	"os"
	"strconv"
)

var (
	appId         int64
	instId        int64
	path          string
	noKeyError    = "you need to pass the private key either with '-k' parameter or by setting 'GITHUB_KEY_PATH' or 'GITHUB_KEY_VALUE' or by passing it via STDIN\n"
	noAppIdError  = "You need to define the App ID via '-a' parameter or 'GITHUB_APP_ID' environment variable\n"
	noInstIdError = "You need to define the Installation ID via '-i' parameter or 'GITHUB_INST_ID' environment variable\n"
)

func ParseParameters() (*GitHubAppConfig, error) {
	var err error
	cfg := &GitHubAppConfig{}

	flag.StringVar(&path, "k", "", "Path to key file for authentication")
	flag.Int64Var(&appId, "a", 0, "App ID to use for authentication")
	flag.Int64Var(&instId, "i", 0, "Installation ID that identifies the APP installation ID on GitHub")
	flag.Parse()

	// Get the values from the environment variables if they are set with parameters
	if appId == 0 {
		appId, err = strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
		if err != nil || appId == 0 {
			return nil, errors.New(noAppIdError)
		}
	}
	if instId == 0 {
		instId, err = strconv.ParseInt(os.Getenv("GITHUB_INST_ID"), 10, 64)
		if err != nil || instId == 0 {
			return nil, errors.New(noInstIdError)
		}
	}

	// Get the key from environment variables
	cfg.PrivateKey = []byte(os.Getenv("GITHUB_KEY_VALUE"))
	if path == "" {
		path = os.Getenv("GITHUB_KEY_PATH")
	}

	// A dash (-) was passed to path as a standard unix option to read from STDIN
	if path == "-" {
		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		cfg.PrivateKey = stdin
	} else if path != "" {
		cfg.PrivateKey, err = os.ReadFile(path)
		if err != nil {
			return nil, err
		}
	}

	if len(cfg.PrivateKey) == 0 {
		return nil, errors.New(noKeyError)
	}

	cfg.ApplicationID = appId
	cfg.InstallationID = instId
	cfg.PrivateKey = cfg.PrivateKey

	return cfg, nil
}

package github

import (
	"flag"
	"os"
	"strconv"
	"testing"
)

var (
	tmpKeyFile string = "/tmp/key.pem"
)

func prepare() {
	instId = 12345
	appId = 9876543
	os.WriteFile(tmpKeyFile, testPrivateKey, 0644)
}

func cleanup() {
	os.Remove(tmpKeyFile)
	os.Args = []string{}
	os.Setenv("GITHUB_APP_ID", "")
	os.Setenv("GITHUB_INST_ID", "")
	os.Setenv("GITHUB_KEY_PATH", "")
	os.Setenv("GITHUB_KEY_VALUE", "")
}

func setArgs(args []string) {
	os.Args = args
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

}

func TestPassedKeyPath(t *testing.T) {
	prepare()
	defer cleanup()

	setArgs([]string{"binary", "-i", strconv.FormatInt(instId, 10), "-a", strconv.FormatInt(appId, 10), "-k", tmpKeyFile})
	cfg, err := ParseParameters()
	if err != nil {
		t.Fatalf("Error validating parameters")
	}
	if cfg.ApplicationID != appId {
		t.Errorf("App ID: expect: %v, got: %v", appId, cfg.ApplicationID)
	}
	if cfg.InstallationID != instId {
		t.Errorf("Installation ID: expect: %v, got: %v", instId, cfg.InstallationID)
	}
	// Read the content of the key file and compare with the parsed private key
	keyContent, err := os.ReadFile(tmpKeyFile)
	if err != nil {
		t.Fatalf("Failed to read key file: %v", err)
	}
	if string(cfg.PrivateKey) != string(keyContent) {
		t.Errorf("Private Key: expected %s, got %s", string(keyContent), string(cfg.PrivateKey))
	}
}

func TestNoAppIDParameter(t *testing.T) {
	prepare()
	defer cleanup()

	setArgs([]string{"binary", "-i", strconv.FormatInt(instId, 10), "-k", tmpKeyFile})
	_, err := ParseParameters()
	if err != nil && err.Error() != noAppIdError {
		t.Errorf("expect: %q, got: %q", noAppIdError, err.Error())
	}
}

func TestNoInstIDParameter(t *testing.T) {
	prepare()
	defer cleanup()

	setArgs([]string{"binary", "-a", strconv.FormatInt(appId, 10), "-k", tmpKeyFile})
	_, err := ParseParameters()
	if err != nil && err.Error() != noInstIdError {
		t.Errorf("expect: %q, got: %q", noInstIdError, err.Error())
	}
}

func TestNoKeyPathParameter(t *testing.T) {
	prepare()
	defer cleanup()

	setArgs([]string{"binary", "-a", strconv.FormatInt(appId, 10), "-i", strconv.FormatInt(instId, 10)})
	_, err := ParseParameters()
	if err != nil && err.Error() != noKeyError {
		t.Errorf("expect: %q, got: %q", noKeyError, err.Error())
	}
}

func TestNoAppIDEnvVar(t *testing.T) {
	prepare()
	defer cleanup()

	setArgs([]string{"binary"})
	os.Setenv("GITHUB_INST_ID", strconv.FormatInt(instId, 10))
	os.Setenv("GITHUB_KEY_PATH", string(tmpKeyFile))
	os.Setenv("GITHUB_KEY_VALUE", string(testPrivateKey))
	_, err := ParseParameters()
	if err != nil && err.Error() != noAppIdError {
		t.Errorf("expect: %q, got: %q", noAppIdError, err.Error())
	}
}

func TestNoInstIDEnvVar(t *testing.T) {
	prepare()
	defer cleanup()

	setArgs([]string{"binary"})
	os.Setenv("GITHUB_APP_ID", strconv.FormatInt(appId, 10))
	os.Setenv("GITHUB_KEY_PATH", string(tmpKeyFile))
	os.Setenv("GITHUB_KEY_VALUE", string(testPrivateKey))
	_, err := ParseParameters()
	if err != nil && err.Error() != noInstIdError {
		t.Errorf("expect: %q, got: %q", noInstIdError, err.Error())
	}
}

func TestNoKeyEnvVar(t *testing.T) {
	prepare()
	defer cleanup()

	setArgs([]string{"binary"})
	os.Setenv("GITHUB_APP_ID", strconv.FormatInt(appId, 10))
	os.Setenv("GITHUB_INST_ID", strconv.FormatInt(instId, 10))
	_, err := ParseParameters()
	if err != nil && err.Error() != noKeyError {
		t.Errorf("expect: %q, got: %q", noKeyError, err.Error())
	}
}

func TestKeyPathEnvVar(t *testing.T) {
	prepare()
	defer cleanup()

	setArgs([]string{"binary"})
	os.Setenv("GITHUB_APP_ID", strconv.FormatInt(appId, 10))
	os.Setenv("GITHUB_INST_ID", strconv.FormatInt(instId, 10))
	os.Setenv("GITHUB_KEY_PATH", string(tmpKeyFile))
	cfg, err := ParseParameters()
	if err != nil {
		t.Fatalf("Uncaught error with environment variables")
	}
	id, _ := strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
	if cfg.ApplicationID != id {
		t.Errorf("GITHUB_APP_ID: expect: %v, got: %v", id, cfg.ApplicationID)
	}
	id, _ = strconv.ParseInt(os.Getenv("GITHUB_INST_ID"), 10, 64)
	if cfg.InstallationID != id {
		t.Errorf("GITHUB_INST_ID: expect: %v, got: %v", id, cfg.InstallationID)
	}

	// Read the content of the key file and compare with the parsed private key
	keyContent, err := os.ReadFile(tmpKeyFile)
	if err != nil {
		t.Fatalf("Failed to read key file: %v", err)
	}
	if string(cfg.PrivateKey) != string(keyContent) {
		t.Errorf("GITHUB_KEY_PATH: expected %s, got %s", string(keyContent), string(cfg.PrivateKey))
	}

}

func TestKeyValueEnvVar(t *testing.T) {
	prepare()
	defer cleanup()

	setArgs([]string{"binary"})
	os.Setenv("GITHUB_APP_ID", strconv.FormatInt(appId, 10))
	os.Setenv("GITHUB_INST_ID", strconv.FormatInt(instId, 10))
	os.Setenv("GITHUB_KEY_VALUE", string(testPrivateKey))
	cfg, err := ParseParameters()
	if err != nil {
		t.Fatalf("Uncaught error with environment variables")
	}
	id, _ := strconv.ParseInt(os.Getenv("GITHUB_APP_ID"), 10, 64)
	if cfg.ApplicationID != id {
		t.Errorf("GITHUB_APP_ID: expect: %v, got: %v", id, cfg.ApplicationID)
	}
	id, _ = strconv.ParseInt(os.Getenv("GITHUB_INST_ID"), 10, 64)
	if cfg.InstallationID != id {
		t.Errorf("GITHUB_INST_ID: expect: %v, got: %v", id, cfg.InstallationID)
	}

	// Read the content of the key file and compare with the parsed private key
	if string(cfg.PrivateKey) != string(testPrivateKey) {
		t.Errorf("GITHUB_KEY_PATH: expected %s, got %s", string(testPrivateKey), string(cfg.PrivateKey))
	}

}

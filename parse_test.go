package github

import (
	"flag"
	"os"
	"strconv"
	"testing"
)

var (
	instID     int64  = 12345
	appID      int64  = 9876543
	tmpKeyFile string = "/tmp/key.pem"
)

func prepare() {
	os.WriteFile(tmpKeyFile, testPrivateKey, 0644)
}

func cleanup() {
	os.Remove(tmpKeyFile)
	os.Args = []string{}
}

func setArgs(args []string) {
	os.Args = args
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

}
func TestPassedKeyPath(t *testing.T) {
	prepare()
	defer cleanup()

	setArgs([]string{"binary", "-i", strconv.FormatInt(instID, 10), "-a", strconv.FormatInt(appID, 10), "-k", tmpKeyFile})
	cfg, err := ParseParameters()
	if err != nil {
		t.Fatalf("Error validating parameters")
	}
	if cfg.ApplicationID != appID {
		t.Errorf("App ID: expect: %v, got: %v", appID, cfg.ApplicationID)
	}
	if cfg.InstallationID != instID {
		t.Errorf("Installation ID: expect: %v, got: %v", instID, cfg.InstallationID)
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

	setArgs([]string{"binary", "-i", strconv.FormatInt(instID, 10), "-k", tmpKeyFile})
	_, err := ParseParameters()
	if err != nil && err.Error() != noAppIdError {
		t.Errorf("expect: %q, got: %q", noAppIdError, err.Error())
	}
}

func TestNoInstIDParameter(t *testing.T) {
	prepare()
	defer cleanup()

	setArgs([]string{"binary", "-a", strconv.FormatInt(appID, 10), "-k", tmpKeyFile})
	_, err := ParseParameters()
	if err != nil && err.Error() != noInstIdError {
		t.Errorf("expect: %q, got: %q", noInstIdError, err.Error())
	}
}

func TestNoKeyPathParameter(t *testing.T) {
	prepare()
	defer cleanup()

	setArgs([]string{"binary", "-a", strconv.FormatInt(appID, 10), "-i", strconv.FormatInt(instID, 10)})
	_, err := ParseParameters()
	if err != nil && err.Error() != noKeyError {
		t.Errorf("expect: %q, got: %q", noKeyError, err.Error())
	}
}

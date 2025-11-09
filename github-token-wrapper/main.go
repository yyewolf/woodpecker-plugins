package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

func main() {
	// Load environment variables
	env, err := loadEnv()
	if err != nil {
		logrus.Fatalf("Error loading environment variables: %v", err)
	}

	if !strings.Contains(env.Plugin.OutputTo, "TOKEN") && !strings.Contains(env.Plugin.OutputTo, "SECRET") {
		logrus.Fatalf("OUTPUT_TO must contain 'TOKEN' or 'SECRET' to avoid accidental overwrites: %s", env.Plugin.OutputTo)
	}

	githubToken, err := GetToken(env.Plugin.MtlsServiceURL, env.Plugin.MtlsClientCert, env.Plugin.MtlsClientKey, env.Plugin.MtlsCACert)
	if err != nil {
		logrus.Fatalf("Error getting GitHub token: %v", err)
	}

	// Run `exec $@` with the token in the environment
	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", env.Plugin.OutputTo, githubToken))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		logrus.Fatalf("Error executing command: %v", err)
	}
}

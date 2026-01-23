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

	// 1. Helm Registry Login
	// Strip scheme if present for login
	registryHost := env.Plugin.Registry
	registryHost = strings.TrimPrefix(registryHost, "http://")
	registryHost = strings.TrimPrefix(registryHost, "https://")
	registryHost = strings.TrimPrefix(registryHost, "oci://")

	if env.Plugin.RegistryUsername != "" && env.Plugin.RegistryPassword != "" {
		logrus.Infof("Logging into registry: %s", registryHost)
		cmd := exec.Command("helm", "registry", "login", registryHost, "--username", env.Plugin.RegistryUsername, "--password-stdin")
		cmd.Stdin = strings.NewReader(env.Plugin.RegistryPassword)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			logrus.Fatalf("Error logging into registry: %v", err)
		}
	} else {
		logrus.Info("Skipping registry login (username or password not set)")
	}

	// 2. Helm Push
	// Construct the OCI URL
	// Schema: oci://<registry>/<path>
	// registryHost is clean, but for push URL we need oci://
	registryPath := strings.TrimPrefix(env.Plugin.RegistryPath, "/")
	ociURL := fmt.Sprintf("oci://%s/%s", registryHost, registryPath)
	logrus.Infof("Pushing chart %s to %s", env.Plugin.ChartPath, ociURL)

	cmd := exec.Command("helm", "push", env.Plugin.ChartPath, ociURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Fatalf("Error pushing chart: %v", err)
	}

	logrus.Infof("Chart pushed successfully")
}

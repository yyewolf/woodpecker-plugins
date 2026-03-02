package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func runHelmCommand(args ...string) {
	cmd := exec.Command("helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Fatalf("Error running helm command (%s): %v", strings.Join(args, " "), err)
	}
}

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
	chartPath := env.Plugin.ChartPath

	if env.Plugin.BuildDependencies {
		if strings.HasSuffix(chartPath, ".tgz") {
			logrus.Fatalf("build_dependencies=true requires chart_path to be a chart directory, got packaged chart: %s", chartPath)
		}

		logrus.Infof("Building chart dependencies for %s", chartPath)
		runHelmCommand("dependency", "build", chartPath)

		logrus.Infof("Packaging chart from %s", chartPath)
		packageArgs := []string{"package", chartPath}
		if env.Plugin.PackageFlags != "" {
			packageArgs = append(packageArgs, strings.Fields(env.Plugin.PackageFlags)...)
		}
		runHelmCommand(packageArgs...)

		chartName := filepath.Base(strings.TrimRight(chartPath, "/"))
		packages, err := filepath.Glob(filepath.Join(".", chartName+"-*.tgz"))
		if err != nil {
			logrus.Fatalf("Error looking up packaged chart: %v", err)
		}
		if len(packages) == 0 {
			logrus.Fatalf("No packaged chart found after helm package for chart %s", chartPath)
		}

		chartPath = packages[len(packages)-1]
		logrus.Infof("Using packaged chart: %s", chartPath)
	}

	// Construct the OCI URL
	// Schema: oci://<registry>/<path>
	// registryHost is clean, but for push URL we need oci://
	registryPath := strings.TrimPrefix(env.Plugin.RegistryPath, "/")
	ociURL := fmt.Sprintf("oci://%s/%s", registryHost, registryPath)
	logrus.Infof("Pushing chart %s to %s", chartPath, ociURL)

	cmd := exec.Command("helm", "push", chartPath, ociURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Fatalf("Error pushing chart: %v", err)
	}

	logrus.Infof("Chart pushed successfully")
}

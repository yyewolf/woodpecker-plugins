package main

import (
	"context"
	"net/http"
	"os"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

func main() {
	env, err := loadEnv()
	if err != nil {
		logrus.Fatalf("Error loading environment variables: %v", err)
	}

	if env.Plugin.GithubPrivateKeyPEMPath == "" && env.Plugin.GithubPrivateKeyPEM == "" {
		logrus.Fatal("Either settings.github_private_key_pem or settings.github_private_key_pem_path must be set")
	}

	keyString := env.Plugin.GithubPrivateKeyPEM
	if env.Plugin.GithubPrivateKeyPEMPath != "" {
		keyString = env.Plugin.GithubPrivateKeyPEMPath
	}

	logrus.Infof("Getting Github Token for repo: %s", env.CI.Repo)
	logrus.Infof("Github App ID: %d", env.Plugin.GithubAppID)
	logrus.Infof("Github Installation ID: %d", env.Plugin.GithubInstallationID)

	// Parse the private key to an *rsa.PrivateKey
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(keyString))
	if err != nil {
		logrus.Fatalf("Error parsing private key: %v", err)
	}

	atr := ghinstallation.NewAppsTransportFromPrivateKey(http.DefaultTransport, env.Plugin.GithubAppID, privateKey)
	itr := ghinstallation.NewFromAppsTransport(atr, env.Plugin.GithubInstallationID)

	token, err := itr.Token(context.TODO())
	if err != nil {
		logrus.Fatalf("Error getting installation token: %v", err)
	}

	if env.Plugin.OutputFile == "" {
		logrus.Infof("No output file specified, sending token to .github_token")
		env.Plugin.OutputFile = ".github_token"
	}

	err = os.WriteFile(env.Plugin.OutputFile, []byte(token), 0600)
	if err != nil {
		logrus.Fatalf("Error writing token to file: %v", err)
	}

	logrus.Infof("Token written to %s successfully", env.Plugin.OutputFile)
}

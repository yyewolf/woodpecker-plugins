package main

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Environnement struct {
	CI struct {
		Repo string `env:"REPO,required"`
	} `envPrefix:"CI_"`

	Plugin struct {
		GithubAppID          int64  `env:"GITHUB_APP_ID,required"`
		GithubInstallationID int64  `env:"GITHUB_INSTALLATION_ID,required"`
		GithubPrivateKeyPEM  string `env:"GITHUB_PRIVATE_KEY_PEM"`

		MtlsCACert     string `env:"MTLS_CA_CERT,required"`
		MtlsServerCert string `env:"MTLS_SERVER_CERT,required"`
		MtlsServerKey  string `env:"MTLS_SERVER_KEY,required"`
	} `envPrefix:"PLUGIN_"`
}

func loadEnv() (Environnement, error) {
	godotenv.Load()
	return env.ParseAs[Environnement]()
}

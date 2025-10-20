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
		GithubAppID             int64  `env:"GITHUB_APP_ID,required"`
		GithubInstallationID    int64  `env:"GITHUB_INSTALLATION_ID,required"`
		GithubPrivateKeyPEM     string `env:"GITHUB_PRIVATE_KEY_PEM"`
		GithubPrivateKeyPEMPath string `env:"GITHUB_PRIVATE_KEY_PEM_PATH,file"`

		OutputFile string `env:"OUTPUT_FILE"`
	} `envPrefix:"PLUGIN_"`
}

func loadEnv() (Environnement, error) {
	godotenv.Load()
	return env.ParseAs[Environnement]()
}

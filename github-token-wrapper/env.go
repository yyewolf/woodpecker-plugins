package main

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Environnement struct {
	Plugin struct {
		MtlsServiceURL string `env:"MTLS_SERVICE_URL,required"`
		MtlsCACert     string `env:"MTLS_CA_CERT,required"`
		MtlsClientCert string `env:"MTLS_CLIENT_CERT,required"`
		MtlsClientKey  string `env:"MTLS_CLIENT_KEY,required"`

		OutputTo string `env:"OUTPUT_TO,required" envDefault:"PLUGIN_GITHUB_TOKEN"`
	} `envPrefix:"PLUGIN_"`
}

func loadEnv() (Environnement, error) {
	godotenv.Load()
	return env.ParseAs[Environnement]()
}

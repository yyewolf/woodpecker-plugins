package main

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Environnement struct {
	Plugin struct {
		Registry         string `env:"REGISTRY" envDefault:"ghcr.io"`
		RegistryUsername string `env:"REGISTRY_USERNAME"`
		RegistryPassword string `env:"REGISTRY_PASSWORD"`
		RegistryPath     string `env:"REGISTRY_PATH,required"`

		BuildDependencies bool   `env:"BUILD_DEPENDENCIES" envDefault:"false"`
		PackageFlags      string `env:"PACKAGE_FLAGS"`

		ChartPath string `env:"CHART_PATH,required"`
	} `envPrefix:"PLUGIN_"`
}

func loadEnv() (Environnement, error) {
	godotenv.Load()
	return env.ParseAs[Environnement]()
}

package main

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Environnement struct {
	CI struct {
		Repo   string `env:"REPO,required"`
		Commit struct {
			PullRequest int `env:"PULL_REQUEST,required"`
		} `envPrefix:"COMMIT_"`
	} `envPrefix:"CI_"`

	Plugin struct {
		GithubToken     string `env:"GITHUB_TOKEN"`
		GithubTokenPath string `env:"GITHUB_TOKEN_PATH,file"`

		Comment       string `env:"COMMENT"`
		CommentPath   string `env:"COMMENT_PATH,file"`
		Repo          string `env:"REPO"`
		PullRequest   int    `env:"PULL_REQUEST"`
		MessageID     string `env:"MESSAGE_ID"`
		UpdateInPlace bool   `env:"UPDATE_IN_PLACE" envDefault:"false"`
	} `envPrefix:"PLUGIN_"`
}

func coalesce[T comparable](values ...T) T {
	for _, v := range values {
		var zero T
		if v != zero {
			return v
		}
	}
	return *new(T)
}

func loadEnv() (Environnement, error) {
	godotenv.Load()
	return env.ParseAs[Environnement]()
}

package main

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Environnement struct {
	CI struct {
		Repo         string `env:"REPO,required"`
		CommitSHA    string `env:"COMMIT_SHA"`
		CommitBranch string `env:"COMMIT_BRANCH"`
	} `envPrefix:"CI_"`

	Plugin struct {
		// GitHub Configuration
		Token      string `env:"TOKEN,required"`
		APIURL     string `env:"API_URL" envDefault:"https://api.github.com"`
		GraphQLURL string `env:"GRAPHQL_URL" envDefault:"https://api.github.com"`
		RepoURL    string `env:"REPO_URL,required"`

		// Release Configuration
		ReleaseType  string `env:"RELEASE_TYPE"`
		Path         string `env:"PATH"`
		TargetBranch string `env:"TARGET_BRANCH"`
		ConfigFile   string `env:"CONFIG_FILE" envDefault:"release-please-config.json"`
		ManifestFile string `env:"MANIFEST_FILE" envDefault:".release-please-manifest.json"`

		// Version Configuration
		ReleaseAs                 string `env:"RELEASE_AS"`
		VersioningStrategy        string `env:"VERSIONING_STRATEGY" envDefault:"default"`
		BumpMinorPreMajor         bool   `env:"BUMP_MINOR_PRE_MAJOR" envDefault:"false"`
		BumpPatchForMinorPreMajor bool   `env:"BUMP_PATCH_FOR_MINOR_PRE_MAJOR" envDefault:"false"`
		PrereleaseType            string `env:"PRERELEASE_TYPE"`
		IncludeVInTags            bool   `env:"INCLUDE_V_IN_TAGS" envDefault:"true"`
		IncludeComponentInTag     bool   `env:"INCLUDE_COMPONENT_IN_TAG" envDefault:"false"`

		// Changelog Configuration
		ChangelogPath     string `env:"CHANGELOG_PATH" envDefault:"CHANGELOG.md"`
		ChangelogType     string `env:"CHANGELOG_TYPE"`
		ChangelogHost     string `env:"CHANGELOG_HOST"`
		ChangelogSections string `env:"CHANGELOG_SECTIONS"`

		// PR Configuration
		SkipGithubPullRequest   bool   `env:"SKIP_GITHUB_PULL_REQUEST" envDefault:"false"`
		SkipLabeling            bool   `env:"SKIP_LABELING" envDefault:"false"`
		Fork                    bool   `env:"FORK" envDefault:"false"`
		DraftPullRequest        bool   `env:"DRAFT_PULL_REQUEST" envDefault:"false"`
		Label                   string `env:"LABEL" envDefault:"autorelease: pending"`
		PullRequestTitlePattern string `env:"PULL_REQUEST_TITLE_PATTERN"`
		PullRequestHeader       string `env:"PULL_REQUEST_HEADER"`
		PullRequestFooter       string `env:"PULL_REQUEST_FOOTER"`

		// Release Configuration
		SkipGithubRelease bool   `env:"SKIP_GITHUB_RELEASE" envDefault:"false"`
		Draft             bool   `env:"DRAFT" envDefault:"false"`
		Prerelease        bool   `env:"PRERELEASE" envDefault:"false"`
		ReleaseLabel      string `env:"RELEASE_LABEL" envDefault:"autorelease: tagged"`

		// Component Configuration
		Component        string `env:"COMPONENT"`
		PackageName      string `env:"PACKAGE_NAME"`
		ComponentNoSpace bool   `env:"COMPONENT_NO_SPACE" envDefault:"false"`

		// Other Configuration
		ProxyServer string   `env:"PROXY_SERVER"`
		ExtraFiles  []string `env:"EXTRA_FILES" envSeparator:","`
		VersionFile string   `env:"VERSION_FILE"`
		Signoff     string   `env:"SIGNOFF"`
		DryRun      bool     `env:"DRY_RUN" envDefault:"false"`
		Snapshot    bool     `env:"SNAPSHOT" envDefault:"false"`

		// Override Configuration
		LatestTagVersion string `env:"LATEST_TAG_VERSION"`
		LatestTagSHA     string `env:"LATEST_TAG_SHA"`
		LatestTagName    string `env:"LATEST_TAG_NAME"`
	} `envPrefix:"PLUGIN_"`
}

func loadEnv() (Environnement, error) {
	godotenv.Load()
	return env.ParseAs[Environnement]()
}

package main

import (
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

	// Set proxy if provided
	if env.Plugin.ProxyServer != "" {
		logrus.Infof("Setting proxy server: %s", env.Plugin.ProxyServer)
		os.Setenv("HTTP_PROXY", "http://"+env.Plugin.ProxyServer)
		os.Setenv("HTTPS_PROXY", "http://"+env.Plugin.ProxyServer)
	}

	// Check if we should skip both PR and release creation
	if env.Plugin.SkipGithubPullRequest && env.Plugin.SkipGithubRelease {
		logrus.Info("Both pull request and release creation are skipped. Nothing to do.")
		return
	}

	// Step 1: Create or update release PR (unless skipped)
	if !env.Plugin.SkipGithubPullRequest {
		logrus.Info("Creating or updating release PR...")
		if err := runReleasePR(env); err != nil {
			logrus.Fatalf("Error creating release PR: %v", err)
		}
	}

	// Step 2: Create GitHub release (unless skipped)
	if !env.Plugin.SkipGithubRelease {
		logrus.Info("Creating GitHub release...")
		if err := runGithubRelease(env); err != nil {
			logrus.Fatalf("Error creating GitHub release: %v", err)
		}
	}

	logrus.Info("Release Please completed successfully!")
}

func runReleasePR(env Environnement) error {
	args := []string{"release-pr"}

	// Add required arguments
	args = append(args, "--repo-url", env.Plugin.RepoURL)

	// Add optional arguments
	if env.Plugin.APIURL != "https://api.github.com" {
		args = append(args, "--api-url", env.Plugin.APIURL)
	}
	if env.Plugin.GraphQLURL != "https://api.github.com" {
		args = append(args, "--graphql-url", env.Plugin.GraphQLURL)
	}
	if env.Plugin.TargetBranch != "" {
		args = append(args, "--target-branch", env.Plugin.TargetBranch)
	}
	if env.Plugin.ReleaseType != "" {
		args = append(args, "--release-type", env.Plugin.ReleaseType)
	}
	if env.Plugin.Path != "" {
		args = append(args, "--path", env.Plugin.Path)
	}
	if env.Plugin.ConfigFile != "release-please-config.json" {
		args = append(args, "--config-file", env.Plugin.ConfigFile)
	}
	if env.Plugin.ManifestFile != ".release-please-manifest.json" {
		args = append(args, "--manifest-file", env.Plugin.ManifestFile)
	}
	if env.Plugin.ReleaseAs != "" {
		args = append(args, "--release-as", env.Plugin.ReleaseAs)
	}
	if env.Plugin.VersioningStrategy != "default" {
		args = append(args, "--versioning-strategy", env.Plugin.VersioningStrategy)
	}
	if env.Plugin.BumpMinorPreMajor {
		args = append(args, "--bump-minor-pre-major")
	}
	if env.Plugin.BumpPatchForMinorPreMajor {
		args = append(args, "--bump-patch-for-minor-pre-major")
	}
	if env.Plugin.PrereleaseType != "" {
		args = append(args, "--prerelease-type", env.Plugin.PrereleaseType)
	}
	if !env.Plugin.IncludeVInTags {
		args = append(args, "--include-v-in-tags", "false")
	}
	if env.Plugin.IncludeComponentInTag {
		args = append(args, "--monorepo-tags")
	}
	if env.Plugin.ChangelogPath != "CHANGELOG.md" {
		args = append(args, "--changelog-path", env.Plugin.ChangelogPath)
	}
	if env.Plugin.ChangelogType != "" {
		args = append(args, "--changelog-type", env.Plugin.ChangelogType)
	}
	if env.Plugin.ChangelogHost != "" {
		args = append(args, "--changelog-host", env.Plugin.ChangelogHost)
	}
	if env.Plugin.ChangelogSections != "" {
		args = append(args, "--changelog-sections", env.Plugin.ChangelogSections)
	}
	if env.Plugin.SkipLabeling {
		args = append(args, "--skip-labeling")
	}
	if env.Plugin.Fork {
		args = append(args, "--fork")
	}
	if env.Plugin.DraftPullRequest {
		args = append(args, "--draft-pull-request")
	}
	if env.Plugin.Label != "autorelease: pending" {
		args = append(args, "--label", env.Plugin.Label)
	}
	if env.Plugin.PullRequestTitlePattern != "" {
		args = append(args, "--pull-request-title-pattern", env.Plugin.PullRequestTitlePattern)
	}
	if env.Plugin.PullRequestHeader != "" {
		args = append(args, "--pull-request-header", env.Plugin.PullRequestHeader)
	}
	if env.Plugin.PullRequestFooter != "" {
		args = append(args, "--pull-request-footer", env.Plugin.PullRequestFooter)
	}
	if env.Plugin.Component != "" {
		args = append(args, "--component", env.Plugin.Component)
	}
	if env.Plugin.PackageName != "" {
		args = append(args, "--package-name", env.Plugin.PackageName)
	}
	if env.Plugin.ComponentNoSpace {
		args = append(args, "--component-no-space")
	}
	if len(env.Plugin.ExtraFiles) > 0 {
		args = append(args, "--extra-files", strings.Join(env.Plugin.ExtraFiles, ","))
	}
	if env.Plugin.VersionFile != "" {
		args = append(args, "--version-file", env.Plugin.VersionFile)
	}
	if env.Plugin.Signoff != "" {
		args = append(args, "--signoff", env.Plugin.Signoff)
	}
	if env.Plugin.DryRun {
		args = append(args, "--dry-run")
	}
	if env.Plugin.Snapshot {
		args = append(args, "--snapshot")
	}
	if env.Plugin.LatestTagVersion != "" {
		args = append(args, "--latest-tag-version", env.Plugin.LatestTagVersion)
	}
	if env.Plugin.LatestTagSHA != "" {
		args = append(args, "--latest-tag-sha", env.Plugin.LatestTagSHA)
	}
	if env.Plugin.LatestTagName != "" {
		args = append(args, "--latest-tag-name", env.Plugin.LatestTagName)
	}

	logrus.Infof("Running: release-please %s", strings.Join(args, " "))

	args = append(args, "--token", env.Plugin.Token)

	cmd := exec.Command("release-please", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runGithubRelease(env Environnement) error {
	args := []string{"github-release"}

	// Add required arguments
	args = append(args, "--repo-url", env.Plugin.RepoURL)

	// Add optional arguments
	if env.Plugin.APIURL != "https://api.github.com" {
		args = append(args, "--api-url", env.Plugin.APIURL)
	}
	if env.Plugin.GraphQLURL != "https://api.github.com" {
		args = append(args, "--graphql-url", env.Plugin.GraphQLURL)
	}
	if env.Plugin.TargetBranch != "" {
		args = append(args, "--target-branch", env.Plugin.TargetBranch)
	}
	if env.Plugin.ReleaseType != "" {
		args = append(args, "--release-type", env.Plugin.ReleaseType)
	}
	if env.Plugin.Path != "" {
		args = append(args, "--path", env.Plugin.Path)
	}
	if env.Plugin.ConfigFile != "release-please-config.json" {
		args = append(args, "--config-file", env.Plugin.ConfigFile)
	}
	if env.Plugin.ManifestFile != ".release-please-manifest.json" {
		args = append(args, "--manifest-file", env.Plugin.ManifestFile)
	}
	if !env.Plugin.IncludeVInTags {
		args = append(args, "--include-v-in-tags", "false")
	}
	if env.Plugin.IncludeComponentInTag {
		args = append(args, "--monorepo-tags")
	}
	if env.Plugin.PullRequestTitlePattern != "" {
		args = append(args, "--pull-request-title-pattern", env.Plugin.PullRequestTitlePattern)
	}
	if env.Plugin.PullRequestHeader != "" {
		args = append(args, "--pull-request-header", env.Plugin.PullRequestHeader)
	}
	if env.Plugin.PullRequestFooter != "" {
		args = append(args, "--pull-request-footer", env.Plugin.PullRequestFooter)
	}
	if env.Plugin.Component != "" {
		args = append(args, "--component", env.Plugin.Component)
	}
	if env.Plugin.PackageName != "" {
		args = append(args, "--package-name", env.Plugin.PackageName)
	}
	if env.Plugin.ComponentNoSpace {
		args = append(args, "--component-no-space")
	}
	if env.Plugin.Draft {
		args = append(args, "--draft")
	}
	if env.Plugin.Prerelease {
		args = append(args, "--prerelease")
	}
	if env.Plugin.Label != "autorelease: pending" {
		args = append(args, "--label", env.Plugin.Label)
	}
	if env.Plugin.ReleaseLabel != "autorelease: tagged" {
		args = append(args, "--release-label", env.Plugin.ReleaseLabel)
	}
	if env.Plugin.DryRun {
		args = append(args, "--dry-run")
	}

	logrus.Infof("Running: release-please %s", strings.Join(args, " "))

	args = append(args, "--token", env.Plugin.Token)

	cmd := exec.Command("release-please", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

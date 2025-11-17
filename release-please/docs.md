# Release Please Plugin

A Woodpecker CI plugin that wraps [Google's release-please](https://github.com/googleapis/release-please) tool for automated releases based on conventional commits.

## Overview

This plugin automates CHANGELOG generation, the creation of GitHub releases, and version bumps for your projects by:

1. Parsing your git history for [Conventional Commit messages](https://www.conventionalcommits.org/)
2. Creating release PRs with updated changelogs and version files
3. Creating GitHub releases when release PRs are merged

## Usage

### Basic Example

```yaml
steps:
  - name: release
    image: yyewolf/woodpecker-plugins:release-please
    settings:
      token:
        from_secret: github_token
      repo_url: https://github.com/owner/repo
      release_type: node
```

### Advanced Example

```yaml
steps:
  - name: release
    image: yyewolf/woodpecker-plugins:release-please
    settings:
      token:
        from_secret: github_token
      repo_url: https://github.com/owner/repo
      release_type: go
      target_branch: main
      path: packages/my-package
      config_file: .release-please-config.json
      manifest_file: .release-please-manifest.json
      changelog_path: CHANGELOG.md
      skip_github_pull_request: false
      skip_github_release: false
      include_component_in_tag: true
      draft_pull_request: false
      fork: false
      dry_run: false
```

## Settings

### Required Settings

| Setting | Description |
|---------|-------------|
| `token` | GitHub token with repo write permissions |
| `repo_url` | GitHub repository URL (e.g., `https://github.com/owner/repo`) |

### GitHub Configuration

| Setting | Default | Description |
|---------|---------|-------------|
| `api_url` | `https://api.github.com` | GitHub API URL |
| `graphql_url` | `https://api.github.com` | GitHub GraphQL URL |

### Release Configuration

| Setting | Default | Description |
|---------|---------|-------------|
| `release_type` | | Type of project (node, go, python, java, etc.) |
| `path` | | Release from path other than root directory |
| `target_branch` | | Branch to open release PRs against |
| `config_file` | `release-please-config.json` | Path to config file |
| `manifest_file` | `.release-please-manifest.json` | Path to manifest file |

### Version Configuration

| Setting | Default | Description |
|---------|---------|-------------|
| `release_as` | | Override the semantically determined release version |
| `versioning_strategy` | `default` | Strategy for version bumping |
| `bump_minor_pre_major` | `false` | Bump minor before first major release |
| `bump_patch_for_minor_pre_major` | `false` | Bump patch instead of minor pre-major |
| `prerelease_type` | | Type of prerelease (e.g., alpha, beta) |
| `include_v_in_tags` | `true` | Include "v" prefix in version tags |
| `include_component_in_tag` | `false` | Include component name in tags (for monorepos) |

### Changelog Configuration

| Setting | Default | Description |
|---------|---------|-------------|
| `changelog_path` | `CHANGELOG.md` | Path to changelog file |
| `changelog_type` | | Type of changelog (default, github) |
| `changelog_host` | | Host for hyperlinks in changelog |
| `changelog_sections` | | Comma-separated scopes to include |

### Pull Request Configuration

| Setting | Default | Description |
|---------|---------|-------------|
| `skip_github_pull_request` | `false` | Skip creating release PRs |
| `skip_labeling` | `false` | Skip labeling PRs |
| `fork` | `false` | Create PR from fork |
| `draft_pull_request` | `false` | Mark PR as draft |
| `label` | `autorelease: pending` | Labels for release PR |
| `pull_request_title_pattern` | | Custom PR title pattern |
| `pull_request_header` | | Custom PR header |
| `pull_request_footer` | | Custom PR footer |

### GitHub Release Configuration

| Setting | Default | Description |
|---------|---------|-------------|
| `skip_github_release` | `false` | Skip creating GitHub releases |
| `draft` | `false` | Mark release as draft |
| `prerelease` | `false` | Mark prerelease versions as prereleases |
| `release_label` | `autorelease: tagged` | Label for completed releases |

### Component Configuration (Monorepos)

| Setting | Default | Description |
|---------|---------|-------------|
| `component` | | Component name for monorepos |
| `package_name` | | Package name being released |
| `component_no_space` | `false` | Disable space before component name |

### Advanced Configuration

| Setting | Default | Description |
|---------|---------|-------------|
| `proxy_server` | | Proxy server (host:port) |
| `extra_files` | | Extra files to consider (comma-separated) |
| `version_file` | | Path to version file |
| `signoff` | | Add signed-off-by line |
| `dry_run` | `false` | Prepare but don't take action |
| `snapshot` | `false` | Generate snapshot/prerelease |

### Override Configuration

| Setting | Default | Description |
|---------|---------|-------------|
| `latest_tag_version` | | Override detected latest tag version |
| `latest_tag_sha` | | Override detected latest tag SHA |
| `latest_tag_name` | | Override detected latest tag name |

## Supported Release Types

- `bazel` - Bazel module with MODULE.bazel
- `dart` - Dart package with pubspec.yaml
- `elixir` - Elixir project with mix.exs
- `go` - Go project
- `helm` - Helm chart with Chart.yaml
- `java` - Java project (generates SNAPSHOT versions)
- `maven` - Maven project (updates pom.xml)
- `node` - Node.js project with package.json
- `expo` - Expo React Native project
- `ocaml` - OCaml project with opam/esy files
- `php` - PHP project with composer.json
- `python` - Python project with pyproject.toml or setup.py
- `ruby` - Ruby project with version.rb
- `rust` - Rust project with Cargo.toml
- `simple` - Simple project with version.txt
- `terraform-module` - Terraform module

## Conventional Commits

This plugin expects your commit messages to follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `fix:` - Bug fixes (patch version bump)
- `feat:` - New features (minor version bump)  
- `feat!:`, `fix!:` - Breaking changes (major version bump)
- `chore:`, `docs:`, `style:` - No version bump
- `BREAKING CHANGE:` in commit body - Major version bump

### Example Commits

```bash
feat: add new API endpoint
fix: resolve memory leak issue
feat!: remove deprecated method
chore: update dependencies
```

## Workflow Examples

### Simple Workflow

```yaml
when:
  branch: main
  event: push

steps:
  - name: release
    image: yyewolf/woodpecker-plugins:release-please
    settings:
      token:
        from_secret: github_token
      repo_url: https://github.com/${CI_REPO}
      release_type: go
```

### Monorepo Workflow

```yaml
when:
  branch: main
  event: push

steps:
  - name: release-package-a
    image: yyewolf/woodpecker-plugins:release-please
    settings:
      token:
        from_secret: github_token
      repo_url: https://github.com/${CI_REPO}
      release_type: node
      path: packages/package-a
      component: package-a
      include_component_in_tag: true

  - name: release-package-b
    image: yyewolf/woodpecker-plugins:release-please
    settings:
      token:
        from_secret: github_token
      repo_url: https://github.com/${CI_REPO}
      release_type: go
      path: packages/package-b
      component: package-b
      include_component_in_tag: true
```

### Conditional Release

```yaml
when:
  branch: main
  event: push

steps:
  - name: create-release-pr
    image: yyewolf/woodpecker-plugins:release-please
    settings:
      token:
        from_secret: github_token
      repo_url: https://github.com/${CI_REPO}
      release_type: node
      skip_github_release: true

  - name: create-github-release
    image: yyewolf/woodpecker-plugins:release-please
    settings:
      token:
        from_secret: github_token
      repo_url: https://github.com/${CI_REPO}
      release_type: node
      skip_github_pull_request: true
    when:
      event: tag
```

## Configuration Files

### release-please-config.json

```json
{
  "release-type": "node",
  "bump-minor-pre-major": true,
  "bump-patch-for-minor-pre-major": true,
  "changelog-sections": [
    {"type": "feat", "section": "Features"},
    {"type": "fix", "section": "Bug Fixes"},
    {"type": "chore", "section": "Miscellaneous", "hidden": true}
  ]
}
```

### .release-please-manifest.json

```json
{
  ".": "1.0.0",
  "packages/frontend": "2.1.0",
  "packages/backend": "1.5.2"
}
```

## Tips

1. **Initial Setup**: Run `release-please bootstrap` locally to create initial configuration files
2. **First Release**: Create an initial release manually or use `release_as` setting
3. **Monorepos**: Use separate paths and components for each package
4. **Testing**: Use `dry_run: true` to test configuration without making changes
5. **Labels**: Release Please uses specific labels to track PR state

## Troubleshooting

- **No Release PR Created**: Ensure you have releasable commits (feat, fix, deps)
- **Wrong Version**: Check if there are existing PRs with `autorelease: pending` labels
- **Permission Errors**: Verify token has necessary repository permissions
- **Configuration Issues**: Use `dry_run: true` to validate settings

For more details, see the [official release-please documentation](https://github.com/googleapis/release-please).
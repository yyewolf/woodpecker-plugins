# Woodpecker Plugins

A collection of Woodpecker CI plugins for GitHub integration.

## Plugins

### [github-app-token](./github-app-token/)

A Woodpecker plugin that generates GitHub App installation tokens for repository access. This plugin is useful when you need to authenticate as a GitHub App instead of using personal access tokens.

**Use cases:**
- Authenticating with GitHub API using GitHub App credentials
- Generating installation tokens for repository operations
- Secure CI/CD workflows with GitHub App permissions

### [github-comment](./github-comment/)

A Woodpecker plugin that creates or updates comments on GitHub pull requests. This plugin allows you to post automated comments during CI/CD workflows.

**Use cases:**
- Posting build status updates to pull requests
- Adding test results or coverage reports
- Creating automated notifications and feedback
- Updating existing comments with new information

## Development

Each plugin is written in Go and can be built as a standalone binary or Docker container.

### Building

To build a specific plugin:

```bash
cd <plugin-directory>
go build -o plugin .
```

### Docker

Each plugin includes a Dockerfile for containerized deployment:

```bash
cd <plugin-directory>
docker build -t <plugin-name> .
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

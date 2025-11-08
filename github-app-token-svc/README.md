# GitHub App Token Plugin

A Woodpecker CI plugin that generates GitHub App installation tokens for repository access.

## Description

This plugin allows you to authenticate with GitHub using a GitHub App instead of personal access tokens. It generates installation tokens that can be used for API calls and repository operations within your CI/CD pipeline.

## Usage

```yaml
steps:
  - name: get-github-token
    image: your-registry/github-app-token
    settings:
      github_app_id: 123456
      github_installation_id: 789012
      github_private_key_pem: |
        -----BEGIN RSA PRIVATE KEY-----
        your-private-key-here
        -----END RSA PRIVATE KEY-----
      output_file: .github_token
```

## Parameters

| Parameter | Required | Description | Default |
|-----------|----------|-------------|---------|
| `github_app_id` | Yes | GitHub App ID | |
| `github_installation_id` | Yes | GitHub App Installation ID | |
| `github_private_key_pem` | No* | GitHub App private key (PEM format) | |
| `github_private_key_pem_path` | No* | Path to file containing GitHub App private key | |
| `output_file` | No | File to write the token to | `.github_token` |

*Either `github_private_key_pem` or `github_private_key_pem_path` must be provided.

## Environment Variables

The plugin uses the following environment variables:

- `CI_REPO` - Repository name (automatically provided by Woodpecker)
- `PLUGIN_GITHUB_APP_ID` - GitHub App ID
- `PLUGIN_GITHUB_INSTALLATION_ID` - GitHub App Installation ID  
- `PLUGIN_GITHUB_PRIVATE_KEY_PEM` - GitHub App private key content
- `PLUGIN_GITHUB_PRIVATE_KEY_PEM_PATH` - Path to private key file
- `PLUGIN_OUTPUT_FILE` - Output file for the token

## Output

The plugin writes the generated GitHub App installation token to the specified output file (default: `.github_token`). This token can then be used by subsequent steps in your pipeline for GitHub API authentication.

## Example Workflow

```yaml
steps:
  - name: generate-token
    image: your-registry/github-app-token
    settings:
      github_app_id: 123456
      github_installation_id: 789012
      github_private_key_pem_path: /secrets/github-app-key.pem
      
  - name: use-token
    image: alpine/git
    commands:
      - export GITHUB_TOKEN=$(cat .github_token)
      - git clone https://x-access-token:$GITHUB_TOKEN@github.com/owner/repo.git
```

## GitHub App Setup

1. Create a GitHub App in your organization or personal account
2. Install the app on the repositories you want to access
3. Note the App ID and Installation ID
4. Generate and download the private key
5. Configure the plugin with these credentials

## Building

```bash
go build -o github-app-token .
```

## Docker

```bash
docker build -t github-app-token .
docker run --rm -e PLUGIN_GITHUB_APP_ID=123456 github-app-token
```
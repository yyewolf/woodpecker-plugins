---
name: GitHub App Token
icon: https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png
author: yyewolf
description: Generate GitHub App installation tokens for repository access
tags: [github, authentication, app, token]
containerImage: ghcr.io/yyewolf/woodpecker-plugins/github-app-token
containerImageUrl: https://github.com/yyewolf/woodpecker-plugins/pkgs/container/woodpecker-plugins%2Fgithub-app-token
url: https://github.com/yyewolf/woodpecker-plugins/tree/main/github-app-token
---

# plugin-github-app-token

Woodpecker plugin to generate GitHub App installation tokens for repository access. This plugin allows you to authenticate with GitHub using a GitHub App instead of personal access tokens, providing better security and granular permissions.

It's goal is to be a singular plugin to allow you to have a secure secret you can use in only this step (without exposing it to the rest of the pipeline) to generate short-lived tokens for accessing GitHub's API.

The generated token can then be used in subsequent steps of your pipeline to perform actions on GitHub, such as commenting on issues, managing pull requests, or triggering workflows.

The plugin is built for the following platforms:
- linux/386
- linux/amd64
- linux/arm
- linux/arm64
- linux/loong64
- linux/mips
- linux/mips64
- linux/mips64le
- linux/mipsle
- linux/ppc64
- linux/ppc64le
- linux/riscv64
- linux/s390x

## Features

- Generate GitHub App installation tokens
- Secure authentication using RSA private keys
- Supports both inline PEM keys and file paths
- Customizable output file location
- Enhanced security over personal access tokens

## Settings

| Settings Name                    | Default         | Description                                                                                          |
| -------------------------------- | --------------- | ---------------------------------------------------------------------------------------------------- |
| `github_app_id`                  | _none_          | **Required.** GitHub App ID (numeric)                                                               |
| `github_installation_id`         | _none_          | **Required.** GitHub App Installation ID (numeric)                                                  |
| `github_private_key_pem`         | _none_          | RSA private key in PEM format (inline). Either this or `github_private_key_pem_path` is required   |
| `github_private_key_pem_path`    | _none_          | Path to file containing RSA private key in PEM format. Either this or `github_private_key_pem` is required |
| `output_file`                    | `.github_token` | File path where the generated token will be written                                                 |

## Examples

### Basic usage with inline PEM key:

```yaml
pipeline:
  get-token:
    image: ghcr.io/yyewolf/woodpecker-plugins/github-app-token
    settings:
      github_app_id: 123456
      github_installation_id: 789012
      github_private_key_pem: |
        -----BEGIN RSA PRIVATE KEY-----
        MIIEpAIBAAKCAQEA1234567890abcdef...
        ...your-private-key-content...
        -----END RSA PRIVATE KEY-----
      output_file: .github_token
```

### Using PEM key from file:

```yaml
pipeline:
  get-token:
    image: ghcr.io/yyewolf/woodpecker-plugins/github-app-token
    settings:
      github_app_id: 123456
      github_installation_id: 789012
      github_private_key_pem_path: /path/to/private-key.pem
```

## Setup Instructions

1. **Create a GitHub App**:
   - Go to your GitHub organization settings
   - Navigate to "Developer settings" → "GitHub Apps"
   - Click "New GitHub App"
   - Configure permissions as needed for your use case

2. **Install the App**:
   - Install the GitHub App on your repositories
   - Note the Installation ID from the installation URL

3. **Generate Private Key**:
   - In your GitHub App settings, generate a private key
   - Download the `.pem` file

4. **Configure Secrets**:
   - Store your App ID, Installation ID, and private key as Woodpecker secrets
   - Use the secrets in your pipeline configuration

## Security Notes

- The generated token has the same permissions as your GitHub App
- Tokens are temporary and automatically expire
- Store private keys securely using Woodpecker secrets
- The output file is created with restricted permissions (0600)
- Never commit private keys to your repository

## Output

The plugin writes the generated token to the specified output file (default: `.github_token`). This token can then be used by subsequent steps in your pipeline for GitHub API operations.

The token file contains only the raw token string and can be read by other tools or scripts in your pipeline.

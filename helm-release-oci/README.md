# GitHub Comment Plugin

A Woodpecker CI plugin that creates or updates comments on GitHub pull requests.

## Description

This plugin allows you to post automated comments to GitHub pull requests during your CI/CD pipeline. It supports both creating new comments and updating existing ones in place, making it perfect for status updates, test results, and other automated feedback.

## Usage

This plugin requires the `github-app-token-svc` service to be running. It uses mTLS to securely authenticate and retrieve a GitHub token, ensuring that only trusted services with valid certificates can perform actions.

You can generate the required PEM certificates (CA, server, and client keys) using the `generate_pem.sh` script provided in the [github-app-token-svc](../github-app-token-svc) directory.

> ⚠️ **Important:** Do not hardcode certificates and keys in your pipeline file. Use Woodpecker CI secrets to store these sensitive values securely.

```yaml
steps:
  - name: comment-pr
    image: your-registry/github-comment
    settings:
      service_url: https://github-app-token-svc:8443/token
      mtls_ca_cert:
        from_secret: mtls_ca_cert
      mtls_client_cert:
        from_secret: mtls_client_cert
      mtls_client_key:
        from_secret: mtls_client_key
      comment: |
        ## Build Status
        ✅ Build completed successfully!
        
        - Tests passed: 142/142
        - Coverage: 85%
      message_id: build-status
      update_in_place: true
```

## Parameters

| Parameter | Required | Description | Default |
|-----------|----------|-------------|---------|
| `service_url` | Yes | URL of the github-app-token-svc | |
| `mtls_ca_cert` | Yes | CA Certificate content | |
| `mtls_client_cert` | Yes | Client Certificate content | |
| `mtls_client_key` | Yes | Client Key content | |
| `comment` | No* | Comment text to post | |
| `comment_path` | No* | Path to file containing comment text | |
| `repo` | No | Repository in format `owner/repo` | Uses `CI_REPO` |
| `pull_request` | No | Pull request number | Uses `CI_COMMIT_PULL_REQUEST` |
| `message_id` | No | Unique identifier for the comment | |
| `update_in_place` | No | Update existing comment instead of creating new | `false` |

*Either `comment` or `comment_path` must be provided.

## Environment Variables

The plugin uses the following environment variables:

- `CI_REPO` - Repository name (automatically provided by Woodpecker)
- `CI_COMMIT_PULL_REQUEST` - Pull request number (automatically provided by Woodpecker)
- `PLUGIN_SERVICE_URL` - URL of the token service
- `PLUGIN_MTLS_CA_CERT` - CA certificate
- `PLUGIN_MTLS_CLIENT_CERT` - Client certificate
- `PLUGIN_MTLS_CLIENT_KEY` - Client private key
- `PLUGIN_COMMENT` - Comment text content
- `PLUGIN_COMMENT_PATH` - Path to comment text file
- `PLUGIN_REPO` - Override repository name
- `PLUGIN_PULL_REQUEST` - Override pull request number
- `PLUGIN_MESSAGE_ID` - Comment identifier for updates
- `PLUGIN_UPDATE_IN_PLACE` - Enable in-place comment updates

## Features

### New Comments

Create a new comment on every run:

```yaml
settings:
  service_url: https://github-app-token-svc:8443/token
  mtls_ca_cert:
    from_secret: mtls_ca_cert
  mtls_client_cert:
    from_secret: mtls_client_cert
  mtls_client_key:
    from_secret: mtls_client_key
  comment: "🚀 Deployment started at $(date)"
```

### Update in Place

Update the same comment on every run using a message ID:

```yaml
settings:
  service_url: https://github-app-token-svc:8443/token
  mtls_ca_cert:
    from_secret: mtls_ca_cert
  mtls_client_cert:
    from_secret: mtls_client_cert
  mtls_client_key:
    from_secret: mtls_client_key
  comment: "📊 Current test results: 95% pass rate"
  message_id: test-results
  update_in_place: true
```

### Comment from File

Load comment content from a file:

```yaml
settings:
  service_url: https://github-app-token-svc:8443/token
  mtls_ca_cert:
    from_secret: mtls_ca_cert
  mtls_client_cert:
    from_secret: mtls_client_cert
  mtls_client_key:
    from_secret: mtls_client_key
  comment_path: ./reports/coverage-summary.md
  message_id: coverage-report
  update_in_place: true
```

## Comment Format

The plugin automatically adds metadata to comments:

```markdown
Your comment content here

<details><summary>Details</summary>
<p>
<sub>
MessageID: #your-message-id  <br>
Written at: 2025-10-20T14:30:00Z  <br>
Written by [woodpecker-plugins/github-comment](https://github.com/yyewolf/woodpecker-plugins/github-comment)  <br>
</p>
</details>
```

## Example Workflows

### Build Status Updates

```yaml
steps:
  - name: build
    image: golang:1.21
    commands:
      - go build ./...
      
  - name: notify-success
    image: your-registry/github-comment
    settings:
      service_url: https://github-app-token-svc:8443/token
      mtls_ca_cert:
        from_secret: mtls_ca_cert
      mtls_client_cert:
        from_secret: mtls_client_cert
      mtls_client_key:
        from_secret: mtls_client_key
      comment: |
        ## ✅ Build Successful
        
        The build completed successfully!
        - Commit: ${CI_COMMIT_SHA}
        - Branch: ${CI_COMMIT_BRANCH}
      message_id: build-status
      update_in_place: true
    when:
      status: [success]
      
  - name: notify-failure
    image: your-registry/github-comment
    settings:
      service_url: https://github-app-token-svc:8443/token
      mtls_ca_cert:
        from_secret: mtls_ca_cert
      mtls_client_cert:
        from_secret: mtls_client_cert
      mtls_client_key:
        from_secret: mtls_client_key
      comment: |
        ## ❌ Build Failed
        
        The build failed. Please check the logs.
        - Commit: ${CI_COMMIT_SHA}
        - Branch: ${CI_COMMIT_BRANCH}
      message_id: build-status
      update_in_place: true
    when:
      status: [failure]
```

### Test Results

```yaml
steps:
  - name: test
    image: golang:1.21
    commands:
      - go test -v ./... | tee test-results.txt
      
  - name: post-results
    image: your-registry/github-comment
    settings:
      service_url: https://github-app-token-svc:8443/token
      mtls_ca_cert:
        from_secret: mtls_ca_cert
      mtls_client_cert:
        from_secret: mtls_client_cert
      mtls_client_key:
        from_secret: mtls_client_key
      comment_path: test-results.txt
      message_id: test-results
      update_in_place: true
```

## Building

### Local Build
```bash
go build -o github-comment .
```

### Docker

#### Single Architecture
```bash
docker build -t github-comment .
docker run --rm \
  -e PLUGIN_SERVICE_URL="https://github-app-token-svc:8443/token" \
  -e PLUGIN_MTLS_CA_CERT="$(cat ca.crt)" \
  -e PLUGIN_MTLS_CLIENT_CERT="$(cat client.crt)" \
  -e PLUGIN_MTLS_CLIENT_KEY="$(cat client.key)" \
  -e PLUGIN_COMMENT="Hello!" \
  github-comment
```

#### Multi-Architecture
```bash
# Create and use a buildx builder
docker buildx create --name multiarch --use

# Build for multiple architectures
docker buildx build \
  --platform linux/amd64,linux/arm64,linux/arm/v7 \
  -t github-comment:multiarch \
  --push .

# Or build and load for local testing (single arch)
docker buildx build \
  --platform linux/amd64 \
  -t github-comment:local \
  --load .
```
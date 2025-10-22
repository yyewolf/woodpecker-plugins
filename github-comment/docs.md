---
name: GitHub Comment
icon: https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png
author: yyewolf
description: Create or update comments on GitHub pull requests
tags: [github, comment, pull-request, automation]
containerImage: ghcr.io/yyewolf/woodpecker-plugins/github-comment
containerImageUrl: https://github.com/yyewolf/woodpecker-plugins/pkgs/container/woodpecker-plugins%2Fgithub-comment
url: https://github.com/yyewolf/woodpecker-plugins/tree/main/github-comment
---

# plugin-github-comment

Woodpecker plugin to create or update comments on GitHub pull requests. Perfect for posting automated feedback, build status, test results, and other CI/CD notifications directly to your pull requests.

To update comments in place (instead of creating new comments each time), use the `message_id` setting to uniquely identify your comment. This allows for cleaner PR discussions by keeping related updates together. Behind the scenes, the plugin uses this ID and puts it in a hidden footer to track and update the comment, if the ID is not unique, you may end up updating the wrong comment.

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

- Create new comments on pull requests
- Update existing comments in place using message IDs
- Support for both inline comments and comment files
- Automatic timestamp and attribution footer
- Flexible repository and PR targeting

## Settings

| Settings Name       | Default                     | Description                                                                                    |
| ------------------- | --------------------------- | ---------------------------------------------------------------------------------------------- |
| `github_token`      | _none_                      | GitHub personal access token or App token. Either this or `github_token_path` is required     |
| `github_token_path` | _none_                      | Path to file containing GitHub token. Either this or `github_token` is required               |
| `comment`           | _none_                      | Comment text (supports Markdown). Either this or `comment_path` is required                   |
| `comment_path`      | _none_                      | Path to file containing comment text. Either this or `comment` is required                    |
| `message_id`        | _none_                      | Unique identifier for the comment. Required when `update_in_place` is true                    |
| `update_in_place`   | `false`                     | If true, updates existing comment with same `message_id` instead of creating new one          |
| `repo`              | _CI_REPO_                   | Repository in format `owner/repo`. Defaults to current repository                             |
| `pull_request`      | _CI_COMMIT_PULL_REQUEST_    | Pull request number. Defaults to current PR                                                   |

## Examples

### Basic comment on current PR:

```yaml
pipeline:
  comment:
    image: ghcr.io/yyewolf/woodpecker-plugins/github-comment
    settings:
      github_token:
        from_secret: github_token
      comment: |
        ## Build Status
        ✅ Build completed successfully!
        
        - Tests: All passed
        - Coverage: 85%
        - Duration: 2m 34s
```

### Update comment in place:

```yaml
pipeline:
  update-status:
    image: ghcr.io/yyewolf/woodpecker-plugins/github-comment
    settings:
      github_token:
        from_secret: github_token
      comment: |
        ## 🔄 Build In Progress
        
        Current step: Running tests...
      message_id: uniquely-identifiable-message-id
      update_in_place: true

  final-status:
    image: ghcr.io/yyewolf/woodpecker-plugins/github-comment
    settings:
      github_token:
        from_secret: github_token
      comment: |
        ## ✅ Build Complete
        
        All tests passed successfully!
      message_id: uniquely-identifiable-message-id
      update_in_place: true
```

### Comment from file:

```yaml
pipeline:
  test-results:
    image: your-test-runner
    commands:
      - run-tests > test-results.md

  comment-results:
    image: ghcr.io/yyewolf/woodpecker-plugins/github-comment
    settings:
      github_token:
        from_secret: github_token
      comment_path: test-results.md
      message_id: uniquely-identifiable-message-id
      update_in_place: true
```

### Comment on different repository:

```yaml
pipeline:
  cross-repo-comment:
    image: ghcr.io/yyewolf/woodpecker-plugins/github-comment
    settings:
      github_token:
        from_secret: github_token
      repo: owner/other-repo
      pull_request: 42
      comment: |
        ## Deployment Status
        
        This change has been deployed to staging environment.
        
        🔗 [View deployment](https://staging.example.com)
```

### Multiple notification strategy:

```yaml
pipeline:
  start-notification:
    image: ghcr.io/yyewolf/woodpecker-plugins/github-comment
    settings:
      github_token:
        from_secret: github_token
      comment: "🚀 Starting build process..."
      message_id: uniquely-identifiable-message-id
      update_in_place: true

  update-progress:
    image: ghcr.io/yyewolf/woodpecker-plugins/github-comment
    settings:
      github_token:
        from_secret: github_token
      comment: "⚡ Running tests and building artifacts..."
      message_id: uniquely-identifiable-message-id
      update_in_place: true

  final-result:
    image: ghcr.io/yyewolf/woodpecker-plugins/github-comment
    settings:
      github_token:
        from_secret: github_token
      comment: |
        ## ✅ Build Complete
        
        - ✅ Tests passed
        - ✅ Build successful  
        - ✅ Artifacts generated
        
        Ready for review! 🎉
      message_id: uniquely-identifiable-message-id
      update_in_place: true
```

## Comment Format

All comments include an automatic footer with:
- Message ID (if provided)
- Timestamp
- Attribution to the plugin

Example footer:
```html
<details><summary>Details</summary>
<p>
<sub>
MessageID: #uniquely-identifiable-message-id<br>
Date: 2025-10-22T14:30:00Z<br>
By <a href="https://github.com/yyewolf/woodpecker-plugins/github-comment">yyewolf/woodpecker-plugins/github-comment</a><br>
</p>
</details>
```

## Authentication

The plugin supports multiple authentication methods:

1. **Personal Access Token**: Create a classic token with `repo` scope
2. **GitHub App Token**: Use with the `github-app-token` plugin for enhanced security
3. **Fine-grained Personal Access Token**: For repository-specific access

### Using with GitHub App Token plugin:

```yaml
pipeline:
  get-token:
    image: yyewolf/woodpecker-plugin-github-app-token
    settings:
      github_app_id:
        from_secret: github_app_id
      github_installation_id:
        from_secret: github_installation_id
      github_private_key_pem:
        from_secret: github_private_key_pem

  comment:
    image: ghcr.io/yyewolf/woodpecker-plugins/github-comment
    settings:
      github_token_path: .github_token
      comment: "Build completed with GitHub App authentication!"
```

## Error Handling

The plugin will fail with descriptive error messages if:
- Required authentication is missing
- Comment content is not provided
- Target repository or PR cannot be found
- GitHub API returns an error
- `update_in_place` is true but `message_id` is not provided

## Best Practices

1. **Use message IDs**: Always use `message_id` for status updates to avoid comment spam
2. **Update in place**: Set `update_in_place: true` for progress updates
3. **Secure tokens**: Store GitHub tokens as Woodpecker secrets
4. **Rich formatting**: Use Markdown for better-formatted comments
5. **Conditional comments**: Use Woodpecker's conditional execution for relevant comments only

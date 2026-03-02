---
name: Helm Release OCI
icon: https://raw.githubusercontent.com/cncf/artwork/main/projects/helm/icon/color/helm-icon-color.png
author: yyewolf
description: Publish Helm charts to OCI registries
tags: [helm, oci, registry, release, kubernetes]
containerImage: ghcr.io/yyewolf/woodpecker-plugins/helm-release-oci
containerImageUrl: https://github.com/yyewolf/woodpecker-plugins/pkgs/container/woodpecker-plugins%2Fhelm-release-oci
url: https://github.com/yyewolf/woodpecker-plugins/tree/main/helm-release-oci
---

# plugin-helm-release-oci

Woodpecker plugin to publish Helm charts to OCI registries using `helm push`. This plugin automates the process of logging into an OCI registry (like ghcr.io, Docker Hub, or private registries) and pushing your Helm charts.

The plugin is built to run on:
- linux/amd64
- linux/arm64

## Features

- Authenticate with OCI registries
- Push Helm charts to OCI registries
- Support for custom registry paths

## Settings

| Settings Name       | Default                     | Description                                                                                    |
| ------------------- | --------------------------- | ---------------------------------------------------------------------------------------------- |
| `registry`          | `ghcr.io`                   | The OCI registry host to connect to                                                            |
| `registry_username` | _none_                      | Username for registry authentication                                                           |
| `registry_password` | _none_                      | Password or token for registry authentication                                                  |
| `registry_path`     | _none_                      | **Required**. The path/namespace in the registry (e.g., `owner/repo` or `library`)             |
| `build_dependencies`| `false`                     | Run `helm dependency build` before publishing. Requires `chart_path` to point to a chart directory |
| `package_flags`     | _none_                      | Extra flags passed to `helm package` (only used when `build_dependencies=true`)                |
| `chart_path`        | _none_                      | **Required**. Path to chart package `.tgz`, or chart directory when `build_dependencies=true` |

## Examples

### Basic usage with GHCR

```yaml
pipeline:
  publish:
    image: ghcr.io/yyewolf/woodpecker-plugins/helm-release-oci
    settings:
      registry: ghcr.io
      registry_username: myuser
      registry_password:
        from_secret: cr_password
      registry_path: my-org/charts
      chart_path: ./dist/my-chart-0.1.0.tgz
```

### Usage with local registry

```yaml
pipeline:
  publish-local:
    image: ghcr.io/yyewolf/woodpecker-plugins/helm-release-oci
    settings:
      registry: localhost:5000
      registry_username: user
      registry_password: password
      registry_path: helm-charts
      chart_path: my-chart-0.1.0.tgz
```

### Unauthenticated Push (if supported by registry)

```yaml
pipeline:
  publish-public:
    image: ghcr.io/yyewolf/woodpecker-plugins/helm-release-oci
    settings:
      registry: public-registry.example.com
      registry_path: public/charts
      chart_path: my-chart-0.1.0.tgz
```

### Build dependencies before push

When `build_dependencies` is enabled, set `chart_path` to the chart directory (not a `.tgz` package). The plugin will run `helm dependency build`, package the chart, and push the generated archive.

```yaml
pipeline:
  publish-with-deps:
    image: ghcr.io/yyewolf/woodpecker-plugins/helm-release-oci
    settings:
      registry: ghcr.io
      registry_username: myuser
      registry_password:
        from_secret: cr_password
      registry_path: my-org/charts
      build_dependencies: true
      package_flags: "--app-version ${CI_COMMIT_TAG} --version ${CI_COMMIT_TAG#v}"
      chart_path: ./charts/my-chart
```

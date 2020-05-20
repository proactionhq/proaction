[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/proactionhq/proaction)

# Proaction

[Proaction](https://proaction.io) is a CLI that recommends and updates GitHub Action Workflows in order to make them more reliable. Certain patterns in Workflows can result in flaky and unreliable output, or can create dependencies that break when external Actions are updated.

The goal of Proaction is to encourage creating workflows that secure, reliable, and will not change unexpectedly.

## Best Practices

Proaction will recommend changes to workflows in order to follow the following best practices:

### 1. Reproducibility

Having reproducible workflows is important in order to ensure that each execution is both reliable and deterministic. A workflow is reproducible when multiple executions of the same workflow using the same commit is guaranteed to produce the exact same result at artifact.

### 2. Secure

Workflows should not use Actions with open CVEs or other security vulnerabilities.

### 3. Updated

Workflows should be able to easily remain updated to use the latest version of an Action. This is needed for security fixes and for performance and feature updates from the Action.

### 4. Maintainability

Workflows should be written to be easy to maintain, minimizing the work needed to follow the other best practices.

## Getting Started

### Install Proaction

To install Proaction, download the latest release from the [Releases](https://github.com/proactionhq/proaction/releases) page or visit the [docs](https://proaction.io/docs/getting-started/installing/) for other options.

### Running With A Workflow File

```shell
$ proaction scan ./path/to/.github/workflows/workflow.yaml
```

### Running In GitHub Actions

Proaction can run as a GitHub Action. To create a workflow using Proaction, [read the Action documentation](https://github.com/proactionhq/proaction/tree/master/action).

## Read More

To read more, visit the [documentation](https://proaction.io). The docs list all of the Proaction checks that are performed and explain the reasons for each.

### GitHub API and Rate Limits

Proaction uses the GitHub API to look up and analyze actions that a workflow references. Unauthenticated requests to the GitHub API are limited to 60 per hour from any single IP address. To increase this and allow Proaction to scan multiple workflows, create a [Personal Access Token](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line) and give this token repo access.

### Data / Telemetry

This project collects usage data to help improve the product. No data will be sent if the `--no-track` flag is provided to the CLI.

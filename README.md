# Proaction

## Installing

## GitHub API and Rate Limits

Proaction uses the GitHub API to look up and analyze actions that a workflow references. Unauthenticated requests to the GitHub API are limited to 60 per hour from any single IP address. To increase this and allow Proaction to scan multiple workflows, create a [Personal Access Token](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line) and give this token repo access.
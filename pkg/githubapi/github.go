package githubapi

import (
	"context"
	"os"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

func NewGitHubClient() *github.Client {
	if os.Getenv("GITHUB_TOKEN") == "" {
		return github.NewClient(nil)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

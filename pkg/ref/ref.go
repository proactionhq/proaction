package ref

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/githubapi"
)

type PossiblyStableTag struct {
	TagName   string
	CommitSHA string
}

type Branch struct {
	BranchName string
	CommitSHA  string
}

// RefToParts takes a uses reference and splits into owner, repo, path and ref
func RefToParts(ref string) (string, string, string, string, error) {
	splitRef := strings.Split(ref, "@")

	if len(splitRef) < 2 {
		return "", "", "", "", errors.New("unsupported reference format")
	}

	repoParts := splitRef[0]
	tag := splitRef[1]

	splitRepoParts := strings.Split(repoParts, "/")
	owner := ""
	repo := ""
	path := ""

	if len(splitRepoParts) > 2 {
		owner = splitRepoParts[0]
		repo = splitRepoParts[1]
		path = strings.Join(splitRepoParts[2:], string(os.PathSeparator))
	} else if len(splitRepoParts) == 2 {
		owner = splitRepoParts[0]
		repo = splitRepoParts[1]
	}

	return owner, repo, path, tag, nil
}

// DetermineGitHubRefType will use the GitHub API to determine if
// a passed ref is a commit, branch or tag
func DetermineGitHubRefType(owner string, repo string, tag string) (*PossiblyStableTag, *Branch, bool, error) {
	githubClient := githubapi.NewGitHubClient()
	tagResponse, githubResponse, err := githubClient.Git.GetRef(context.Background(), owner, repo, fmt.Sprintf("tags/%s", tag))
	if err != nil {
		if githubResponse.Response.StatusCode != 404 {
			return nil, nil, false, errors.Wrap(err, "failed to get tag ref")
		}
	}

	if tagResponse != nil {
		getTagResponse, _, err := githubClient.Git.GetTag(context.Background(), owner, repo, tagResponse.Object.GetSHA())
		if err != nil {
			return nil, nil, false, errors.Wrap(err, "failed to get commit sha from tag")
		}

		return &PossiblyStableTag{
			TagName:   tag,
			CommitSHA: getTagResponse.Object.GetSHA()[0:7],
		}, nil, false, nil
	}

	branchResponse, githubResponse, err := githubClient.Git.GetRef(context.Background(), owner, repo, fmt.Sprintf("heads/%s", tag))
	if err != nil {
		if githubResponse.Response.StatusCode != 404 {
			return nil, nil, false, errors.Wrap(err, "failed to get head ref")
		}
	}

	if branchResponse != nil {
		return nil, &Branch{
			BranchName: tag,
			CommitSHA:  branchResponse.Object.GetSHA()[0:7],
		}, false, nil
	}

	commitResponse, githubResponse, err := githubClient.Git.GetRef(context.Background(), owner, repo, tag)
	if err != nil {
		if githubResponse.Response.StatusCode != 404 {
			return nil, nil, false, errors.Wrap(err, "failed to get commit ref")
		}
	}

	if commitResponse != nil {
		return nil, nil, true, nil
	}

	return nil, nil, false, nil
}

package ref

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v28/github"
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
func DetermineGitHubRefType(owner string, repo string, unknownRef string) (*PossiblyStableTag, *Branch, bool, error) {
	githubClient := githubapi.NewGitHubClient()
	tagResponse, githubResponse, err := githubClient.Git.GetRef(context.Background(), owner, repo, fmt.Sprintf("tags/%s", unknownRef))
	if err != nil {
		if githubResponse.Response.StatusCode != 404 && err.Error() != "multiple matches found for this ref" {
			return nil, nil, false, errors.Wrapf(err, "failed to get tag ref for owner %s, repo %s, tag %s", owner, repo, unknownRef)
		}
	}

	if tagResponse != nil {
		return &PossiblyStableTag{
			TagName:   unknownRef,
			CommitSHA: tagResponse.Object.GetSHA()[0:7],
		}, nil, false, nil
	}

	branchResponse, githubResponse, err := githubClient.Git.GetRef(context.Background(), owner, repo, fmt.Sprintf("heads/%s", unknownRef))
	if err != nil {
		if githubResponse.Response.StatusCode != 404 {
			return nil, nil, false, errors.Wrap(err, "failed to get head ref")
		}
	}

	if branchResponse != nil {
		return nil, &Branch{
			BranchName: unknownRef,
			CommitSHA:  branchResponse.Object.GetSHA()[0:7],
		}, false, nil
	}

	commitResponse, githubResponse, err := githubClient.Repositories.GetCommit(context.Background(), owner, repo, unknownRef)
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

func IsSHAInRepo(owner string, repo string, commitSHA string) (bool, error) {
	githubClient := githubapi.NewGitHubClient()

	// the github api returns the same result if a commit
	// is present only in a fork... so we need to be
	// a little creative here

	opts := github.ListOptions{}
	listBranchesResponse, _, err := githubClient.Repositories.ListBranches(context.Background(), owner, repo, &opts)
	if err != nil {
		return false, errors.Wrap(err, "failed to list branches for repo")
	}

	for _, branch := range listBranchesResponse {
		compareResponse, _, err := githubClient.Repositories.CompareCommits(context.Background(), owner, repo, branch.GetName(), commitSHA)
		if err != nil {
			return false, errors.Wrap(err, "failed to compare commits")
		}

		if compareResponse.GetStatus() == "behind" || compareResponse.GetStatus() == "identical" {
			return true, nil
		}

		// ahead or diverged are not in the branch
	}

	return false, nil
}

func ListTagsInRepo(owner string, repo string) ([]string, error) {
	githubClient := githubapi.NewGitHubClient()

	opts := github.ListOptions{}
	listTagsResponse, _, err := githubClient.Repositories.ListTags(context.Background(), owner, repo, &opts)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to list tags for repo")
	}

	tags := []string{}
	for _, tag := range listTagsResponse {
		tags = append(tags, tag.GetName())
	}

	return tags, nil
}

func CreateRefString(owner string, repo string, path string, ref string) string {
	if path == "" {
		return fmt.Sprintf("%s/%s@%s", owner, repo, ref)
	}

	return fmt.Sprintf("%s/%s/%s@%s", owner, repo, path, ref)
}

package collect

import (
	"context"
	"fmt"

	"github.com/google/go-github/v28/github"
	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/collect/types"
	"github.com/proactionhq/proaction/pkg/githubapi"
	"github.com/proactionhq/proaction/pkg/logger"
	"github.com/proactionhq/proaction/pkg/ref"
	"go.uber.org/zap"
)

func parseGitHubRef(workflowInfo types.WorkflowInfo, collectors []string) (*types.Output, error) {
	logger.Debug("parseGitHubRef",
		zap.String("input", workflowInfo.LineContent),
		zap.Strings("collectors", collectors))

	output := types.Output{}
	repoOutput := types.RepoOutput{
		WorkflowInfo: workflowInfo,
	}

	owner, repo, path, unknownRef, err := ref.RefToParts(workflowInfo.LineContent)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse ref")
	}

	for _, collector := range collectors {
		if collector == "repo.info" {
			err := retrieveRepoInfo(owner, repo, path, &repoOutput)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get repo info")
			}
		} else if collector == "repo.ref" {
			err := retrieveRefInfo(owner, repo, unknownRef, &repoOutput)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get ref info")
			}
		} else if collector == "repo.branches" {
			err := retreiveBranches(owner, repo, &repoOutput)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get branches")
			}
		} else if collector == "repo.commits" {
			err := retreiveCommits(owner, repo, &repoOutput)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get commits")
			}
		} else if collector == "repo.tags" {
			err := retreiveTags(owner, repo, &repoOutput)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get tags")
			}
		} else {
			return nil, errors.Errorf("unknown collector %q", collector)
		}
	}

	output.Repos = []*types.RepoOutput{&repoOutput}

	return &output, nil
}

func retrieveRepoInfo(owner string, repo string, path string, output *types.RepoOutput) error {
	githubClient := githubapi.NewGitHubClient()

	githubRepo, _, err := githubClient.Repositories.Get(context.Background(), owner, repo)
	if err != nil {
		return errors.Wrap(err, "failed to get github repo")
	}

	output.ID = githubRepo.GetID()
	output.Owner = githubRepo.GetOwner().GetLogin()
	output.Repo = githubRepo.GetName()
	output.IsArchived = githubRepo.GetArchived()
	output.IsPublic = true // TODO the version of the client doesn't have this field
	output.DefaultBranch = githubRepo.GetDefaultBranch()
	output.Forks = []string{} // TODO
	output.IsFork = githubRepo.GetFork()
	output.Head = "" // TODO
	if githubRepo.GetParent() != nil {
		output.Parent = githubRepo.GetParent().GetFullName()
	}

	if path != "" {
		output.Path = fmt.Sprintf("/%s", path)
	}

	return nil
}

func retrieveRefInfo(owner string, repo string, unknownRef string, output *types.RepoOutput) error {
	possibleTag, maybeBranch, isCommit, err := ref.DetermineGitHubRefType(owner, repo, unknownRef)
	if err != nil {
		return errors.Wrap(err, "failed to determine ref type")
	}

	if possibleTag != nil {
		output.Ref = unknownRef
		output.RefType = "tag"
	} else if maybeBranch != nil {
		output.Ref = unknownRef
		output.RefType = "branch"
	} else if isCommit {
		output.Ref = unknownRef
		output.RefType = "commit"
	}

	return nil
}

func retreiveBranches(owner string, repo string, output *types.RepoOutput) error {
	githubClient := githubapi.NewGitHubClient()

	githubBranches, _, err := githubClient.Repositories.ListBranches(
		context.Background(),
		owner,
		repo,
		&github.ListOptions{},
	)
	if err != nil {
		return errors.Wrap(err, "failed to list branches")
	}

	output.Branches = []types.BranchOutput{}
	for _, branch := range githubBranches {
		output.Branches = append(output.Branches, types.BranchOutput{
			Name: branch.GetName(),
			Head: branch.GetCommit().GetSHA(),
		})
	}

	return nil
}

func retreiveTags(owner string, repo string, output *types.RepoOutput) error {
	githubClient := githubapi.NewGitHubClient()

	githubTags, _, err := githubClient.Repositories.ListTags(
		context.Background(),
		owner,
		repo,
		&github.ListOptions{},
	)
	if err != nil {
		return errors.Wrap(err, "failed to list tags")
	}

	output.Tags = []types.TagOutput{}
	for _, tag := range githubTags {
		output.Tags = append(output.Tags, types.TagOutput{
			Name: tag.GetName(),
			Head: tag.GetCommit().GetSHA(),
		})
	}

	return nil
}

func retreiveCommits(owner string, repo string, output *types.RepoOutput) error {
	githubClient := githubapi.NewGitHubClient()

	githubCommits, _, err := githubClient.Repositories.ListCommits(
		context.Background(),
		owner,
		repo,
		&github.CommitsListOptions{},
	)
	if err != nil {
		return errors.Wrap(err, "failed to list commits")
	}

	output.Commits = []string{}
	for _, commit := range githubCommits {
		output.Commits = append(output.Commits, commit.GetSHA())
	}

	return nil
}

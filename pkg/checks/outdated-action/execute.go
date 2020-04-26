package outdatedaction

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/githubapi"
	"github.com/proactionhq/proaction/pkg/issue"
	"github.com/proactionhq/proaction/pkg/ref"
	workflowtypes "github.com/proactionhq/proaction/pkg/workflow/types"
)

func executeOutdatedActionCheckForWorkflow(parsedWorkflow *workflowtypes.GitHubWorkflow) ([]*issue.Issue, error) {
	issues := []*issue.Issue{}

	for jobName, job := range parsedWorkflow.Jobs {
		for stepIdx, step := range job.Steps {
			if step.Uses == "" {
				continue
			}

			// ignore docker uses
			if strings.HasPrefix(step.Uses, "docker://") {
				continue
			}

			owner, repo, path, tag, err := ref.RefToParts(step.Uses)
			if err != nil {
				return nil, errors.Wrap(err, "failed to parse ref")
			}

			_, _, isCommit, err := ref.DetermineGitHubRefType(owner, repo, tag)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get ref type")
			}

			if !isCommit {
				continue
			}

			if len(tag) > 7 {
				tag = tag[0:7]
			}

			// Get the latest commit from the repo
			githubClient := githubapi.NewGitHubClient()
			getRepoResponse, _, err := githubClient.Repositories.Get(context.Background(), owner, repo)
			if err != nil {
				return nil, errors.Wrap(err, "failed to get repo")
			}

			getBranchResponse, _, err := githubClient.Repositories.GetBranch(context.Background(), owner, repo, getRepoResponse.GetDefaultBranch())
			if err != nil {
				return nil, errors.Wrap(err, "failed to get branch")
			}

			latestCommit := getBranchResponse.GetCommit().GetSHA()[0:7]

			if tag == latestCommit {
				continue
			}

			stableRef := ""
			if path == "" {
				stableRef = fmt.Sprintf("%s/%s@%s", owner, repo, latestCommit)
			} else {
				stableRef = fmt.Sprintf("%s/%s/%s@%s", owner, repo, path, latestCommit)
			}

			message := mustGetIssueMessage(parsedWorkflow.Name, jobName, step)

			i := issue.Issue{
				CheckType: CheckName,
				JobName:   jobName,
				StepIdx:   stepIdx,
				CheckData: map[string]interface{}{
					"originalGitHubRef":   step.Uses,
					"remediatedGitHubRef": stableRef,
				},
				Message:      message,
				CanRemediate: true,
			}

			issues = append(issues, &i)
		}
	}

	return issues, nil
}

func mustGetIssueMessage(workflowName string, jobName string, step *workflowtypes.Step) string {
	return fmt.Sprintf("The job named %q in the %q workflow is referencing an outdated commit from %q.", jobName, workflowName, step.Uses)
}

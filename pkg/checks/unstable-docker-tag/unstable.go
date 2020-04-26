package unstabledockertag

import (
	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/issue"
	workflowtypes "github.com/proactionhq/proaction/pkg/workflow/types"
)

var (
	CheckName = "unstable-docker-tag"
)

func Run(originalWorkflowContent string, parsedWorkflow *workflowtypes.GitHubWorkflow) ([]*issue.Issue, error) {
	issues, err := executeUnstableTagCheckForWorkflow(parsedWorkflow)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute unstable tag check")
	}

	for _, i := range issues {
		err := remediateIssue(parsedWorkflow, i)
		if err != nil {
			return nil, errors.Wrap(err, "failed to remediate workflow")
		}
	}

	return issues, nil
}

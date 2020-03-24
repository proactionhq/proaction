package unstabledockertag

import (
	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/issue"
	workflowtypes "github.com/proactionhq/proaction/pkg/workflow/types"
)

var (
	CheckName = "unstable-docker-tag"
)

func Run(originalWorkflowContent string, parsedWorkflow *workflowtypes.GitHubWorkflow) (string, []*issue.Issue, error) {
	issues, err := executeUnstableTagCheckForWorkflow(parsedWorkflow)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to execute unstable tag check")
	}

	lastRemediateWorkflowContent := originalWorkflowContent
	for _, i := range issues {
		afterWorkflow, err := remediateWorkflow(parsedWorkflow, lastRemediateWorkflowContent, i)
		if err != nil {
			return "", nil, errors.Wrap(err, "failed to remediate workflow")
		}

		lastRemediateWorkflowContent = afterWorkflow
	}

	return lastRemediateWorkflowContent, issues, nil
}

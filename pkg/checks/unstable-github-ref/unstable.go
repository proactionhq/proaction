package unstablegithubref

import (
	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/issue"
	workflowtypes "github.com/proactionhq/proaction/pkg/workflow/types"
)

var (
	CheckName = "unstable-github-ref"
)

func Run(originalWorkflowContent string, parsedWorkflow *workflowtypes.GitHubWorkflow) (string, []*issue.Issue, error) {
	issues, err := executeUnstableRefCheckForWorkflow(parsedWorkflow)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to execute unstable ref check")
	}

	lastRemediateWorkflowContent := originalWorkflowContent
	for _, i := range issues {
		afterWorkflow, err := remediateWorkflow(lastRemediateWorkflowContent, i)
		if err != nil {
			return "", nil, errors.Wrap(err, "failed to remediate workflow")
		}

		lastRemediateWorkflowContent = afterWorkflow
	}

	return lastRemediateWorkflowContent, issues, nil
}

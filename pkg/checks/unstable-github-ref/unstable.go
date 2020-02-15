package unstablegithubref

import (
	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/issue"
	"github.com/proactionhq/proaction/pkg/workflow"
)

var (
	CheckName = "unstable-github-ref"
)

func Run(originalWorkflowContent string, parsedWorkflow *workflow.ParsedWorkflow) (string, []*issue.Issue, error) {
	issues, err := executeUnstableRefCheckForWorkflow(parsedWorkflow)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to execute unstable ref check")
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

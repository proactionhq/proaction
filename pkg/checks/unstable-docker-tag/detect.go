package unstabledockertag

import (
	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/issue"
	workflowtypes "github.com/proactionhq/proaction/pkg/workflow/types"
)

var (
	CheckName = "unstable-docker-tag"
)

// DetectIssues will analyze the parsedWorkflow and return a list of issues
func DetectIssues(parsedWorkflow workflowtypes.GitHubWorkflow) ([]*issue.Issue, error) {
	issues, err := executeUnstableTagCheckForWorkflow(parsedWorkflow)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute unstable tag check")
	}

	return issues, nil
}

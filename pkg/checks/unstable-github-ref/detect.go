package unstablegithubref

import (
	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/issue"
	workflowtypes "github.com/proactionhq/proaction/pkg/workflow/types"
)

var (
	CheckName = "unstable-github-ref"
)

// DetectIssues will analyze the parsedWorkflow and return a list of issues
func DetectIssues(parsedWorkflow workflowtypes.GitHubWorkflow) ([]*issue.Issue, error) {
	issues, err := executeUnstableRefCheckForWorkflow(parsedWorkflow)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute unstable ref check")
	}

	return issues, nil
}

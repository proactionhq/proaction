package unforkaction

import (
	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/issue"
	workflowtypes "github.com/proactionhq/proaction/pkg/workflow/types"
)

var (
	CheckName = "unfork-action"
)

func DetectIssues(parsedWorkflow workflowtypes.GitHubWorkflow) ([]*issue.Issue, error) {
	issues, err := executeUnforkActionCheckForWorkflow(parsedWorkflow)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute unfork action check")
	}

	return issues, nil
}

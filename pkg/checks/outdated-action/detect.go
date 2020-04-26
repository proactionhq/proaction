package outdatedaction

import (
	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/issue"
	workflowtypes "github.com/proactionhq/proaction/pkg/workflow/types"
)

var (
	CheckName = "outdated-action"
)

func DetectIssues(parsedWorkflow workflowtypes.GitHubWorkflow) ([]*issue.Issue, error) {
	issues, err := executeOutdatedActionCheckForWorkflow(parsedWorkflow)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute outdated action check")
	}

	return issues, nil
}

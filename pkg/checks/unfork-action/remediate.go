package unforkaction

import (
	"errors"
	"strings"

	"github.com/proactionhq/proaction/pkg/issue"
	workflowtypes "github.com/proactionhq/proaction/pkg/workflow/types"
)

func remediateIssue(parsedWorkflow *workflowtypes.GitHubWorkflow, i *issue.Issue) error {
	job, ok := parsedWorkflow.Jobs[i.JobName]
	if !ok {
		return errors.New("job not found in workflow")
	}

	step := job.Steps[i.StepIdx]

	step.Uses = strings.ReplaceAll(step.Uses, i.CheckData["originalGitHubRef"].(string), i.CheckData["remediatedGitHubRef"].(string))
	job.Steps[i.StepIdx] = step

	return nil
}

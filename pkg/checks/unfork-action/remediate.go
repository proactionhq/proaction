package unforkaction

import (
	"strings"

	"github.com/proactionhq/proaction/pkg/issue"
	"github.com/proactionhq/proaction/pkg/workflow"
)

func remediateWorkflow(parsedWorkflow *workflow.ParsedWorkflow, beforeWorkflowContent string, i *issue.Issue) (string, error) {
	// we do a string replace here because... we don't want to lose comments and rework
	// too much of the yaml

	updatedContent := strings.ReplaceAll(beforeWorkflowContent, i.CheckData["originalGitHubRef"].(string), i.CheckData["remediatedGitHubRef"].(string))
	return updatedContent, nil
}

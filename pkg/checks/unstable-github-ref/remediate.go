package unstablegithubref

import (
	"strings"

	"github.com/proactionhq/proaction/pkg/issue"
)

func remediateWorkflow(beforeWorkflowContent string, i *issue.Issue) (string, error) {
	// we do a string replace here because... we don't want to lose comments and rework
	// too much of the yaml

	updatedContent := strings.ReplaceAll(beforeWorkflowContent, i.CheckData["originalGitHubRef"].(string), i.CheckData["remediatedGitHubRef"].(string))
	return updatedContent, nil
}

package unstablegithubref

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/issue"
	"github.com/proactionhq/proaction/pkg/workflow"
)

func remediateWorkflow(parsedWorkflow *workflow.ParsedWorkflow, beforeWorkflowContent string, i *issue.Issue) (string, error) {
	// todo stop checking this a second time, we already checked it
	owner, repo, _, tag, err := refToParts(i.CheckData["githubRef"].(string))
	if err != nil {
		return "", errors.Wrap(err, "failed to split reference")
	}

	possiblyStableTag, maybeBranch, _, err := determineGitHubRefType(0, owner, repo, tag)
	if err != nil {
		return "", errors.Wrap(err, "failed to get ref type")
	}

	updatedRef := ""
	if maybeBranch != nil {
		updatedRef = fmt.Sprintf("%s/%s@%s", owner, repo, maybeBranch.CommitSHA)
	} else if possiblyStableTag != nil {
		updatedRef = fmt.Sprintf("%s/%s@%s", owner, repo, possiblyStableTag.CommitSHA)
	}

	// we do a string replace here because... we don't want to lose comments and rework
	// too much of the yaml
	updatedContent := strings.ReplaceAll(beforeWorkflowContent, i.CheckData["githubRef"].(string), updatedRef)
	return updatedContent, nil
}

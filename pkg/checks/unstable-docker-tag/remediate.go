package unstabledockertag

import (
	"strings"

	"github.com/proactionhq/proaction/pkg/issue"
)

func RemediateIssue(content string, i *issue.Issue) (string, error) {
	lines := strings.Split(content, "\n")

	lines[i.LineNumber-1] = strings.ReplaceAll(
		lines[i.LineNumber-1],
		i.CheckData["originalGitHubRef"].(string),
		i.CheckData["remediatedGitHubRef"].(string),
	)

	return strings.Join(lines, "\n"), nil
}

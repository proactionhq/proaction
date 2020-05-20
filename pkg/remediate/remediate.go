package remediate

import (
	"strings"

	"github.com/proactionhq/proaction/pkg/remediate/types"
)

func Remediate(input string, lineNumber int, originalContent string, remediatedContent string) (*types.Remediation, error) {
	lines := strings.Split(input, "\n")

	before := lines[lineNumber-1]
	lines[lineNumber-1] = strings.ReplaceAll(
		lines[lineNumber-1],
		originalContent,
		remediatedContent,
	)
	after := lines[lineNumber-1]

	output := input
	if before != after {
		output = strings.Join(lines, "\n")
	}

	remediation := types.Remediation{
		StartLine:         lineNumber,
		OriginalContent:   originalContent,
		RemediatedContent: remediatedContent,

		BeforeWorkflow: originalContent,
		AfterWorkflow:  output,
		WasRemediated:  before != after,
	}

	return &remediation, nil
}

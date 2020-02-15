package scanner

import (
	"fmt"

	"github.com/pkg/errors"
	unstablegithubref "github.com/proactionhq/proaction/pkg/checks/unstable-github-ref"
	"github.com/proactionhq/proaction/pkg/issue"
	"github.com/proactionhq/proaction/pkg/workflow"
)

type Scanner struct {
	OriginalContent   string
	RemediatedContent string
	Issues            []*issue.Issue
}

func NewScanner() *Scanner {
	return &Scanner{
		Issues: []*issue.Issue{},
	}
}

func (s *Scanner) ScanWorkflow(parsedWorkflow *workflow.ParsedWorkflow) error {
	updatedContent, issues, err := unstablegithubref.Run(s.OriginalContent, parsedWorkflow)
	if err != nil {
		return errors.Wrap(err, "failed to run unstable github ref check")
	}

	s.RemediatedContent = updatedContent
	s.Issues = append(s.Issues, issues...)

	return nil
}

func (s Scanner) GetOutput() string {
	output := ""
	for _, i := range s.Issues {
		output = fmt.Sprintf("%s--->%s\n", output, i.Message)
	}

	return output
}

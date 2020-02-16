package scanner

import (
	"fmt"

	"github.com/pkg/errors"
	outdatedaction "github.com/proactionhq/proaction/pkg/checks/outdated-action"
	unforkaction "github.com/proactionhq/proaction/pkg/checks/unfork-action"
	unstabledockertag "github.com/proactionhq/proaction/pkg/checks/unstable-docker-tag"
	unstablegithubref "github.com/proactionhq/proaction/pkg/checks/unstable-github-ref"
	"github.com/proactionhq/proaction/pkg/issue"
	"github.com/proactionhq/proaction/pkg/workflow"
)

type Scanner struct {
	OriginalContent   string
	RemediatedContent string
	Issues            []*issue.Issue
	EnabledChecks     []string
}

func NewScanner() *Scanner {
	return &Scanner{
		Issues:        []*issue.Issue{},
		EnabledChecks: []string{},
	}
}

func (s *Scanner) EnableAllChecks() {
	s.EnabledChecks = []string{
		"unstable-github-ref",
		"unstable-docker-tag",
		"outdated-action",
		"unfork-action",
	}
}

func (s *Scanner) ScanWorkflow(parsedWorkflow *workflow.ParsedWorkflow) error {
	for _, check := range s.EnabledChecks {
		if check == "unstable-github-ref" {
			updatedContent, issues, err := unstablegithubref.Run(s.OriginalContent, parsedWorkflow)
			if err != nil {
				return errors.Wrap(err, "failed to run unstable unstable-github ref check")
			}

			s.RemediatedContent = updatedContent
			s.Issues = append(s.Issues, issues...)
		} else if check == "unstable-docker-tag" {
			updatedContent, issues, err := unstabledockertag.Run(s.OriginalContent, parsedWorkflow)
			if err != nil {
				return errors.Wrap(err, "failed to run unstable unstable-docker-tag check")
			}

			s.RemediatedContent = updatedContent
			s.Issues = append(s.Issues, issues...)
		} else if check == "outdated-action" {
			updatedContent, issues, err := outdatedaction.Run(s.OriginalContent, parsedWorkflow)
			if err != nil {
				return errors.Wrap(err, "failed to run unstable outdated-action check")
			}

			s.RemediatedContent = updatedContent
			s.Issues = append(s.Issues, issues...)
		} else if check == "unfork-action" {
			updatedContent, issues, err := unforkaction.Run(s.OriginalContent, parsedWorkflow)
			if err != nil {
				return errors.Wrap(err, "failed to run unstable unfork-action check")
			}

			s.RemediatedContent = updatedContent
			s.Issues = append(s.Issues, issues...)
		}
	}

	return nil
}

func (s Scanner) GetOutput() string {
	output := ""
	for _, i := range s.Issues {
		output = fmt.Sprintf("%s--->%s\n", output, i.Message)
	}

	return output
}

package scanner

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"
	outdatedaction "github.com/proactionhq/proaction/pkg/checks/outdated-action"
	unforkaction "github.com/proactionhq/proaction/pkg/checks/unfork-action"
	unstabledockertag "github.com/proactionhq/proaction/pkg/checks/unstable-docker-tag"
	unstablegithubref "github.com/proactionhq/proaction/pkg/checks/unstable-github-ref"
	"github.com/proactionhq/proaction/pkg/issue"
	workflowtypes "github.com/proactionhq/proaction/pkg/workflow/types"
	"gopkg.in/yaml.v3"
)

type ScannerStatus int

const (
	ScannerStatusPending   ScannerStatus = iota
	ScannerStatusRunning   ScannerStatus = iota
	ScannerStatusCompleted ScannerStatus = iota
)

type ScannerProgress struct {
	Steps      []string
	StepStatus map[string]ScannerStatus
}

type ScannerFunc func() ScannerProgress

type Scanner struct {
	OriginalContent   string
	RemediatedContent string
	Issues            []*issue.Issue
	EnabledChecks     []string
	ParsedWorkflow    *workflowtypes.GitHubWorkflow
	JobNames          []string
	Progress          map[string]ScannerFunc
}

func NewScanner(content string) (*Scanner, error) {
	parsedWorkflow := workflowtypes.GitHubWorkflow{}
	if err := yaml.Unmarshal([]byte(content), &parsedWorkflow); err != nil {
		return nil, errors.Wrap(err, "failed to parse content")
	}

	jobNames := []string{}
	for jobName := range parsedWorkflow.Jobs {
		jobNames = append(jobNames, jobName)
	}

	return &Scanner{
		OriginalContent: content,
		Issues:          []*issue.Issue{},
		EnabledChecks:   []string{},
		ParsedWorkflow:  &parsedWorkflow,
		JobNames:        jobNames,
	}, nil
}

func (s *Scanner) EnableChecks(checks []string) {
	s.EnabledChecks = checks
	s.initProgress()
}

func (s *Scanner) EnableAllChecks() {
	s.EnabledChecks = []string{
		"unfork-action",
		"unstable-docker-tag",
		"unstable-github-ref",
		"outdated-action",
	}
	s.initProgress()
}

func (s *Scanner) initProgress() {
	s.Progress = map[string]ScannerFunc{}

	for _, enabledCheck := range s.EnabledChecks {
		if enabledCheck == "unstable-github-ref" {
			s.Progress[enabledCheck] = s.progressUnstableGitHubRef
		}
	}
}

func (s Scanner) progressUnstableGitHubRef() ScannerProgress {
	return ScannerProgress{
		Steps: []string{
			"a",
			"b",
			"c",
			"d",
		},
		StepStatus: map[string]ScannerStatus{
			"a": ScannerStatusCompleted,
			"b": ScannerStatusRunning,
			"c": ScannerStatusPending,
			"d": ScannerStatusPending,
		},
	}
}

func (s *Scanner) ScanWorkflow() error {
	sort.Sort(byPriority(s.EnabledChecks))

	for _, check := range s.EnabledChecks {
		// unmarshal from content each time so that each step can build on the last
		// this is important because if an issue changes the line count in the workflow
		// doing this will allow all remediation to still target the correct lines

		parsedWorkflow := workflowtypes.GitHubWorkflow{}
		if err := yaml.Unmarshal([]byte(s.getContent()), &parsedWorkflow); err != nil {
			return errors.Wrap(err, "failed to parse workflow")
		}

		if check == "unstable-github-ref" {
			issues, err := unstablegithubref.DetectIssues(parsedWorkflow)
			if err != nil {
				return errors.Wrap(err, "failed to run unstable unstable-github ref check")
			}

			s.Issues = append(s.Issues, issues...)

			for _, i := range issues {
				updated, err := unstablegithubref.RemediateIssue(s.getContent(), i)
				if err != nil {
					return errors.Wrap(err, "failed to apply remediation")
				}
				s.RemediatedContent = updated
			}
		} else if check == "unstable-docker-tag" {
			issues, err := unstabledockertag.DetectIssues(parsedWorkflow)
			if err != nil {
				return errors.Wrap(err, "failed to run unstable unstable-docker-tag check")
			}

			s.Issues = append(s.Issues, issues...)

			for _, i := range issues {
				updated, err := unstabledockertag.RemediateIssue(s.getContent(), i)
				if err != nil {
					return errors.Wrap(err, "failed to apply remediation")
				}
				s.RemediatedContent = updated
			}
		} else if check == "outdated-action" {
			issues, err := outdatedaction.DetectIssues(parsedWorkflow)
			if err != nil {
				return errors.Wrap(err, "failed to run unstable outdated-action check")
			}

			s.Issues = append(s.Issues, issues...)

			for _, i := range issues {
				updated, err := outdatedaction.RemediateIssue(s.getContent(), i)
				if err != nil {
					return errors.Wrap(err, "failed to apply remediation")
				}
				s.RemediatedContent = updated
			}
		} else if check == "unfork-action" {
			issues, err := unforkaction.DetectIssues(parsedWorkflow)
			if err != nil {
				return errors.Wrap(err, "failed to run unstable unfork-action check")
			}

			s.Issues = append(s.Issues, issues...)

			for _, i := range issues {
				updated, err := unforkaction.RemediateIssue(s.getContent(), i)
				if err != nil {
					return errors.Wrap(err, "failed to apply remediation")
				}
				s.RemediatedContent = updated
			}
		}
	}

	return nil
}

func applyRemediation(content string, i issue.Issue) (string, error) {
	return content, nil
}

func (s Scanner) getContent() string {
	if s.RemediatedContent != "" {
		return s.RemediatedContent
	}

	return s.OriginalContent
}

func (s Scanner) GetOutput() string {
	output := ""
	for _, i := range s.Issues {
		output = fmt.Sprintf("%s* %s\n", output, i.Message)
	}

	return output
}

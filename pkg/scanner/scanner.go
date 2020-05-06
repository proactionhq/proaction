package scanner

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/checks"
	checktypes "github.com/proactionhq/proaction/pkg/checks/types"
	"github.com/proactionhq/proaction/pkg/collect"
	"github.com/proactionhq/proaction/pkg/collect/types"
	collecttypes "github.com/proactionhq/proaction/pkg/collect/types"
	"github.com/proactionhq/proaction/pkg/evaluate"
	evaluatetypes "github.com/proactionhq/proaction/pkg/evaluate/types"
	"github.com/proactionhq/proaction/pkg/issue"
	"github.com/proactionhq/proaction/pkg/remediate"
	remediatetypes "github.com/proactionhq/proaction/pkg/remediate/types"
	workflowtypes "github.com/proactionhq/proaction/pkg/workflow/types"
	"gopkg.in/yaml.v3"
)

type Scanner struct {
	Filename          string
	OriginalContent   []byte
	RemediatedContent []byte
	Issues            []*issue.Issue
	EnabledChecks     []*checktypes.Check
	ParsedWorkflow    *workflowtypes.GitHubWorkflow
	JobNames          []string

	Results []evaluatetypes.EvaluateResult
}

func NewScanner(filename string, content []byte) (*Scanner, error) {
	parsedWorkflow := workflowtypes.GitHubWorkflow{}
	if err := yaml.Unmarshal(content, &parsedWorkflow); err != nil {
		return nil, errors.Wrap(err, "failed to parse content")
	}

	jobNames := []string{}
	for jobName := range parsedWorkflow.Jobs {
		jobNames = append(jobNames, jobName)
	}

	return &Scanner{
		Filename:        filename,
		OriginalContent: content,
		Issues:          []*issue.Issue{},
		EnabledChecks:   []*checktypes.Check{},
		ParsedWorkflow:  &parsedWorkflow,
		JobNames:        jobNames,
	}, nil
}

func (s *Scanner) EnableChecks(checks []*checktypes.Check) {
	s.EnabledChecks = checks
}

func (s *Scanner) EnableAllChecks() {
	s.EnabledChecks = []*checktypes.Check{
		checks.UnstableGitHubRef(),
		checks.Outdated(),
		checks.Unfork(),
	}
}

func (s *Scanner) ScanWorkflow() error {
	// build collectors
	collectors := []collecttypes.Collector{}
	for _, enabledCheck := range s.EnabledChecks {
		for _, checkCollector := range enabledCheck.Collectors {
			mergedCollectors, err := collect.MergeCollectors(checkCollector, collectors)
			if err != nil {
				return errors.Wrap(err, "failed to merge collectors")
			}

			collectors = mergedCollectors
		}
	}

	// execute collect phase
	output := types.Output{
		Repos: []*types.RepoOutput{},
	}

	for _, collector := range collectors {
		// TODO is filename right here?
		collectorOutput, err := collect.Collect(collector, s.Filename, s.OriginalContent)
		if err != nil {
			return errors.Wrap(err, "failed to collect collector")
		}

		output.Repos = append(output.Repos, collectorOutput.Repos...)
		output.Recommendations = collectorOutput.Recommendations
	}
	marshalledCollectors, err := json.Marshal(output)
	if err != nil {
		return errors.Wrap(err, "failed to marshal combined collectors")
	}

	// evaluation
	evaluateResults := []evaluatetypes.EvaluateResult{}
	for _, enabledCheck := range s.EnabledChecks {
		for _, checkEvaluator := range enabledCheck.Evaluators {
			evaluateOutput, err := evaluate.Evaluate(checkEvaluator, marshalledCollectors)
			if err != nil {
				return errors.Wrap(err, "failed to evaluate")
			}

			evaluateResults = append(evaluateResults, evaluateOutput...)
		}
	}

	reducedResults, err := evaluate.ReduceResults(evaluateResults)
	if err != nil {
		return errors.Wrap(err, "failed to reduce results")
	}

	s.Results = reducedResults

	// remediate
	remediateResults := []remediatetypes.Remediation{}
	s.RemediatedContent = s.OriginalContent
	for _, evaluateResult := range s.Results {
		// TODO this will not work if any remediation changes the line numbers
		// we need a better solution

		remediation, err := remediate.Remediate(
			string(s.RemediatedContent),
			evaluateResult.OriginalLineNumber,
			evaluateResult.OriginalContent,
			evaluateResult.RemediatedContent,
		)
		if err != nil {
			return errors.Wrap(err, "failed to apply remediation")
		}

		remediateResults = append(remediateResults, *remediation)

		if remediation.WasRemediated {
			s.RemediatedContent = []byte(remediation.AfterWorkflow)
		}
	}

	return nil
}

func applyRemediation(content string, i issue.Issue) (string, error) {
	return content, nil
}

func (s Scanner) getContent() []byte {
	if s.RemediatedContent != nil {
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

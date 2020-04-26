package unstabledockertag

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/docker"
	"github.com/proactionhq/proaction/pkg/issue"
	progresstypes "github.com/proactionhq/proaction/pkg/progress/types"
	workflowtypes "github.com/proactionhq/proaction/pkg/workflow/types"
)

type UnstableReason int

const (
	UnknownReason      UnstableReason = iota
	IsLatestTag        UnstableReason = iota
	HasUnstableHistory UnstableReason = iota
)

var (
	CheckName = "unstable-docker-tag"
)

// DetectIssues will analyze the parsedWorkflow and return a list of issues
func DetectIssues(parsedWorkflow workflowtypes.GitHubWorkflow, setProgressFunc progresstypes.SetProgressFunc) ([]*issue.Issue, error) {
	issues := []*issue.Issue{}

	for jobName, job := range parsedWorkflow.Jobs {
		setProgressFunc(jobName, true, false)
		for stepIdx, step := range job.Steps {
			if step.Uses.Value == "" {
				continue
			}

			if !strings.HasPrefix(step.Uses.Value, "docker://") {
				continue
			}

			isStable, unstableReason, err := isImageTagStable(step.Uses.Value)
			if err != nil {
				return nil, errors.Wrap(err, "failed to check is image name stable")
			}

			if isStable {
				continue
			}

			message := mustGetIssueMessage(parsedWorkflow.Name, jobName, unstableReason, step)

			i := issue.Issue{
				CheckType:  CheckName,
				JobName:    jobName,
				StepIdx:    stepIdx,
				LineNumber: step.Uses.Line,

				CheckData: map[string]interface{}{
					"unstableReason": unstableReason,
					"originalTag":    "",
					"redmediatedTag": "",
				},
				Message:      message,
				CanRemediate: true,
			}

			issues = append(issues, &i)
		}
		setProgressFunc(jobName, false, true)
	}

	return issues, nil
}

func mustGetIssueMessage(workflowName string, jobName string, unstableReason UnstableReason, step *workflowtypes.Step) string {
	switch unstableReason {
	case IsLatestTag:
		return fmt.Sprintf("The job named %q in the %q workflow is referencing an action that uses the latest tag of the %q docker image. The latest is likely to change", jobName, workflowName, step.Uses.Value)
	case HasUnstableHistory:
		return "has unstable history"
	}

	return ""
}

func isImageTagStable(imageURI string) (bool, UnstableReason, error) {
	_, _, tag, err := docker.ParseImageName(strings.TrimPrefix(imageURI, "docker://"))
	if err != nil {
		return false, UnknownReason, errors.Wrap(err, "failed to parse image")
	}

	if tag == "latest" {
		return false, IsLatestTag, nil
	}

	return true, UnknownReason, nil
}

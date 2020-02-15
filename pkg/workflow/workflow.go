package workflow

import (
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type ParsedWorkflow struct {
	Name      string                       `yaml:"name,omitempty"`
	Jobs      map[string]ParsedWorkflowJob `yaml:"jobs"`
	Platforms []string
}

type ParsedWorklowStep struct {
	Uses string `yaml:"uses,omitempty"`
}

type ParsedWorkflowJob struct {
	Name   string              `yaml:"name,omitempty"`
	RunsOn string              `yaml:"runs-on,omitempty"`
	Steps  []ParsedWorklowStep `yaml:"steps,omitempty"`
}

func Parse(workflowContent []byte) (*ParsedWorkflow, error) {
	parsedWorkflow := ParsedWorkflow{}
	if err := yaml.Unmarshal(workflowContent, &parsedWorkflow); err != nil {
		return nil, errors.Wrap(err, "failed to parse workflow")
	}

	platforms := []string{}
	for _, job := range parsedWorkflow.Jobs {
		if job.RunsOn == "" {
			platforms = append(platforms, "default")
		} else if strings.HasPrefix(job.RunsOn, "macos-") {
			platforms = append(platforms, "macos")
		} else if strings.HasPrefix(job.RunsOn, "windows-") {
			platforms = append(platforms, "windows")
		} else if strings.HasPrefix(job.RunsOn, "ubuntu") {
			platforms = append(platforms, "linux")
		} else {
			platforms = append(platforms, "other")
		}
	}

	parsedWorkflow.Platforms = platforms

	return &parsedWorkflow, nil
}

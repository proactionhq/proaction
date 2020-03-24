package workflow

import "github.com/proactionhq/proaction/pkg/workflow/types"

func Platforms(workflow *types.GitHubWorkflow) []string {
	platforms := []string{}

	for _, job := range workflow.Jobs {
		if job.RunsOn == nil {
			continue
		}

		if job.RunsOn.Type == types.String {
			platforms = append(platforms, *job.RunsOn.StrVal)
			continue
		}

		for _, runsOn := range job.RunsOn.ListVal {
			platforms = append(platforms, runsOn)
		}
	}

	return platforms
}

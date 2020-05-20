package checks

import (
	"github.com/proactionhq/proaction/pkg/checks/types"
	collecttypes "github.com/proactionhq/proaction/pkg/collect/types"
	evaluatetypes "github.com/proactionhq/proaction/pkg/evaluate/types"
)

func UnstableGitHubRef() *types.Check {
	return &types.Check{
		Collectors: []collecttypes.Collector{
			{
				Name:   "uses",
				Path:   "jobs[*].steps[*].uses",
				Parser: "githubref",
				Collectors: []string{
					"repo.info",
					"repo.ref",
					"repo.branches",
					"repo.tags",
					"repo.commits",
				},
			},
		},
		Evaluators: []evaluatetypes.Evaluator{
			{
				Name: "unstable",
				Rego: UnstableGithubRef_Policy,
			},
		},
	}
}

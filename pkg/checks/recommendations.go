package checks

import (
	"github.com/proactionhq/proaction/pkg/checks/types"
	collecttypes "github.com/proactionhq/proaction/pkg/collect/types"
	evaluatetypes "github.com/proactionhq/proaction/pkg/evaluate/types"
)

func Recommendations() *types.Check {
	return &types.Check{
		Collectors: []collecttypes.Collector{
			{
				Name:   "uses",
				Path:   "jobs[*].steps[*].uses",
				Parser: "githubref",
				Collectors: []string{
					"repo.info",
					"repo.ref",
				},
			},
		},
		Evaluators: []evaluatetypes.Evaluator{
			{
				Name: "recommendation",
				Rego: Recommendations_Policy,
			},
		},
	}
}

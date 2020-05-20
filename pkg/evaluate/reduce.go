package evaluate

import (
	"fmt"
	"sort"

	"github.com/proactionhq/proaction/pkg/evaluate/types"
)

// ReduceResults takes a set of results that may contain multiple
// recommendations for the same line, and reduces them to
// ensure that all of the results can be applied
func ReduceResults(results []types.EvaluateResult) ([]types.EvaluateResult, error) {
	mapped := map[string][]types.EvaluateResult{}

	// map results by workflow / owner / repo / line number
	for _, result := range results {
		key := fmt.Sprintf("%s-%s/%s:%d", result.Workflow, result.Owner, result.Repo, result.OriginalLineNumber)

		existingMapped, ok := mapped[key]
		if ok {
			existingMapped = append(existingMapped, result)
		} else {
			existingMapped = []types.EvaluateResult{result}
		}

		mapped[key] = existingMapped
	}

	reduced := []types.EvaluateResult{}

	for _, toReduce := range mapped {
		sort.Sort(types.ByPriority(toReduce))

		if len(toReduce) == 0 {
			continue
		}

		if toReduce[0].Reason == "isStaticRecommendation" {
			reduced = append(reduced, toReduce[0])
			continue
		}

		// if there's no static, just sort by ref type and take the top
		sort.Sort(types.ByRefType(toReduce))
		reduced = append(reduced, toReduce[0])
	}

	return reduced, nil
}

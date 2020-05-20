package collect

import "github.com/proactionhq/proaction/pkg/collect/types"

// MergeCollectors will merge toMerge into the array of collectors,
// without duplication
func MergeCollectors(toMerge types.Collector, collectors []types.Collector) ([]types.Collector, error) {
	mergedCollectors := []types.Collector{}

	wasMerged := false
	for _, collector := range collectors {
		if collector.Equals(toMerge) {
			collector.Merge(toMerge)
			wasMerged = true
		}

		mergedCollectors = append(mergedCollectors, collector)
	}

	if !wasMerged {
		mergedCollectors = append(mergedCollectors, toMerge)
	}

	return mergedCollectors, nil
}

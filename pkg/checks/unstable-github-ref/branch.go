package unstablegithubref

import (
	"fmt"

	"github.com/pkg/errors"
)

func getUpdatedRefFromBranch(owner string, repo string, path string, branchName string, commitSHA string) (string, error) {
	if recommendationActionRefs == nil {
		if err := loadRecommendationsFromWeb(); err != nil {
			return "", errors.Wrap(err, "failed to load recommendations from web")
		}
	}

	// if the action is in the recommendations
	k := fmt.Sprintf("%s/%s", owner, repo)
	if path != "" {
		k = fmt.Sprintf("%s/%s", k, path)
	}
	recommend, ok := recommendationActionRefs[k]
	if ok {
		if recommend.RecommendedRefType == "tag" {
			if path == "" {
				return fmt.Sprintf("%s/%s@%s", owner, repo, recommend.RecommendedRefs[0]), nil
			} else {
				return fmt.Sprintf("%s/%s/%s@%s", owner, repo, path, recommend.RecommendedRefs[0]), nil
			}
		}
	}

	if path == "" {
		return fmt.Sprintf("%s/%s@%s", owner, repo, commitSHA), nil
	}

	return fmt.Sprintf("%s/%s/%s@%s", owner, repo, path, commitSHA), nil
}

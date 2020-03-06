package unstablegithubref

import (
	"fmt"
	"sort"

	semver "github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/ref"
)

// getUpdatedRefFromTag will return the recommendation when a tag is found
// this is not part of the oudated check.
// it should make every effort to give a more specific version of the same commit
// we don't have all commit shas here, so we imply some best practices
// and make some assumptions about how tags are named
func getUpdatedRefFromTag(owner string, repo string, path string, tagName string, commitSHA string, allTags []string) (string, error) {
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
		// the repo (owner/repo/path) was found in the recommends list
		if recommend.RecommendedRefType == "tag" {
			// the owner recommends tags
			for _, recommendedTag := range recommend.RecommendedRefs {
				if recommendedTag == tagName {
					// exact match, return it
					return ref.CreateRefString(owner, repo, path, tagName), nil
				}
			}

			// the exact tag was not found in the recommendations
			// so let's sort by semver and return the highest
			parsedVersions := []*semver.Version{}
			for _, allTag := range allTags {
				v, err := semver.NewVersion(allTag)
				if err == nil {
					parsedVersions = append(parsedVersions, v)
				}
			}
			if len(parsedVersions) > 0 {
				// there are some semvers so we can pick the highest
				sort.Sort(semver.Collection(parsedVersions))
				highestVersion := parsedVersions[len(parsedVersions)-1]
				return ref.CreateRefString(owner, repo, path, highestVersion.Original()), nil
			}

			// ok so here, we recommend tags, tags aren't semver and we don't have a match
			// this is sorted, so we will recommend the top tag in the list
			if len(allTags) > 0 {
				return ref.CreateRefString(owner, repo, path, allTags[0]), nil
			}
		}
	}

	// if the version parses as semver
	v, err := semver.NewVersion(tagName)
	if err == nil {
		parsedVersions := []*semver.Version{}
		for _, allTag := range allTags {
			v, err := semver.NewVersion(allTag)
			if err == nil {
				parsedVersions = append(parsedVersions, v)
			}
		}
		if len(parsedVersions) > 0 {
			// there are some semvers so we can pick the highest
			sort.Sort(semver.Collection(parsedVersions))

			// Find the most specific version of the tag
			mostSpecificVersion := v
			for _, parsedVersion := range parsedVersions {
				if parsedVersion.GreaterThan(v) {
					if parsedVersion.Major() == v.Major() {
						if parsedVersion.GreaterThan(mostSpecificVersion) {
							mostSpecificVersion = parsedVersion
						}
					}
				}
			}

			return ref.CreateRefString(owner, repo, path, mostSpecificVersion.Original()), nil
		}

	}

	// when all else fails, return the commit sha
	return ref.CreateRefString(owner, repo, path, commitSHA), nil
}

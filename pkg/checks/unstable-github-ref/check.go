package unstablegithubref

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/ref"
)

func isGitHubRefStable(githubRef string) (bool, UnstableReason, string, error) {
	// relative paths are very stable
	if strings.HasPrefix(githubRef, ".") {
		return true, IsStable, githubRef, nil
	}

	// if there's no @ sign, then it's unstable
	if !strings.Contains(githubRef, "@") {
		return true, NoSpecifiedVersion, githubRef, nil
	}

	owner, repo, path, tag, err := ref.RefToParts(githubRef)
	if err != nil {
		return false, UnknownReason, "", errors.Wrap(err, "failed to split ref")
	}

	possiblyStableTag, maybeBranch, isCommit, err := ref.DetermineGitHubRefType(owner, repo, tag)
	if err != nil {
		return false, UnknownReason, "", errors.Wrap(err, "failed to get ref type")
	}

	updatedRef := ""
	if maybeBranch != nil {
		if path == "" {
			updatedRef = fmt.Sprintf("%s/%s@%s", owner, repo, maybeBranch.CommitSHA)
		} else {
			updatedRef = fmt.Sprintf("%s/%s/%s@%s", owner, repo, path, maybeBranch.CommitSHA)
		}
	} else if possiblyStableTag != nil {
		allTags, err := ref.ListTagsInRepo(owner, repo)
		if err != nil {
			return false, UnknownReason, "", errors.Wrap(err, "failed to list all tags")
		}

		updatedRefFromTag, err := getUpdatedRefFromTag(owner, repo, path, possiblyStableTag.TagName, possiblyStableTag.CommitSHA, allTags)
		if err != nil {
			return false, UnknownReason, "", errors.Wrap(err, "failed to get updated tag ref")
		}
		if updatedRefFromTag == githubRef {
			return true, IsStable, githubRef, nil
		}

		return false, NotRecommendedTag, updatedRefFromTag, nil
	}

	if tag == "master" {
		return false, IsMaster, updatedRef, nil
	}

	// first check out cache, see if we know anything about this combination
	isCached, cachedIsStable, cachedUnstableReason, err := isGitHubRefStableInCache(owner, repo, tag)
	if err != nil {
		return false, UnknownReason, "", errors.Wrap(err, "failed to check cache")
	}

	if isCached {
		return cachedIsStable, cachedUnstableReason, updatedRef, nil
	}

	isStable := false
	unstableReason := UnknownReason

	if maybeBranch != nil {
		isStable = false
		unstableReason = IsBranch
	} else if isCommit {
		isStable = true
		unstableReason = IsStable
	} else if possiblyStableTag != nil {
		if err := cacheGitHubTagHistory(owner, repo, tag, possiblyStableTag.CommitSHA); err != nil {
			return false, UnknownReason, updatedRef, errors.Wrap(err, "failed to cache tag history")
		}
		hasUnstableHistory, err := doesTagHaveUnstableHistory(owner, repo, tag)
		if err != nil {
			return false, UnknownReason, updatedRef, errors.Wrap(err, "failed to check if tag has unstable history")
		}

		if hasUnstableHistory {
			isStable = false
			unstableReason = HasUnstableHistory
		} else {
			// now we are in a gray area.
			// it's probably stable, but that's by convention
			isStable = true
			unstableReason = IsStable
		}
	} else {
		// whoa, this isn't a valid tag
		isStable = false
		unstableReason = TagNotFound
	}

	// add to the cache
	cacheDuration := time.Hour * 24 * 30
	if possiblyStableTag != nil {
		cacheDuration = time.Hour * 24 * 3
	}
	if err := cacheGitHubRefStable(owner, repo, tag, isStable, unstableReason, cacheDuration); err != nil {
		// dont fail, but this will chew through rate limits
		fmt.Printf("err")
	}

	return isStable, unstableReason, updatedRef, nil
}

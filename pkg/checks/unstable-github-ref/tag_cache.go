package unstablegithubref

import (
	"time"

	"github.com/spf13/viper"
)

// cacheGitHubTagHistory will cache a hash of the repo along with the tag and commit sha
// on the remote erver. this will be skipped if the --no-track option is set.
// this is cached so that proaction is able to suggest which tags are stable and which
// have proven unstable histories
func cacheGitHubTagHistory(owner string, repo string, tag string, sha string) error {
	v := viper.GetViper()
	if v.GetBool("no-track") {
		return nil
	}

	return nil
}

// doesTagHaveUnstableHistory uses a remote api, and looks up the hashed owner/repo with the
// specified tag, to determine if there have been multiple commits on this single tag, which
// indicastes that it's unstable. this function is not used if --no-track is set
func doesTagHaveUnstableHistory(owner string, repo string, tag string) (bool, error) {
	v := viper.GetViper()
	if v.GetBool("no-track") {
		return false, nil
	}

	return false, nil
}

func isGitHubRefStableInCache(owner string, repo string, tag string) (bool, bool, UnstableReason, error) {
	return false, false, UnknownReason, nil
}

func cacheGitHubRefStable(owner string, repo string, tag string, isStable bool, unstableReason UnstableReason, cacheDuration time.Duration) error {
	return nil
}

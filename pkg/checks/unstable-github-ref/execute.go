package unstablegithubref

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/issue"
	"github.com/proactionhq/proaction/pkg/ref"
	"github.com/proactionhq/proaction/pkg/workflow"
)

type UnstableReason int

const (
	IsStable           UnstableReason = iota
	UnknownReason      UnstableReason = iota
	UnsupportedRef     UnstableReason = iota
	NoSpecifiedVersion UnstableReason = iota
	IsMaster           UnstableReason = iota
	IsBranch           UnstableReason = iota
	HasUnstableHistory UnstableReason = iota
	TagNotFound        UnstableReason = iota
)

func executeUnstableRefCheckForWorkflow(parsedWorkflow *workflow.ParsedWorkflow) ([]*issue.Issue, error) {
	issues := []*issue.Issue{}

	for jobName, job := range parsedWorkflow.Jobs {
		for _, step := range job.Steps {
			if step.Uses == "" {
				continue
			}

			// ignore docker uses
			if strings.HasPrefix(step.Uses, "docker://") {
				continue
			}

			isStable, unstableReason, stableRef, err := isGitHubRefStable(step.Uses)
			if err != nil {
				return nil, errors.Wrap(err, "failed to check is github ref stable")
			}

			if isStable {
				continue
			}

			message := mustGetIssueMessage(parsedWorkflow.Name, jobName, unstableReason, step)

			i := issue.Issue{
				CheckType: CheckName,
				CheckData: map[string]interface{}{
					"jobName":             jobName,
					"unstableReason":      unstableReason,
					"originalGitHubRef":   step.Uses,
					"remediatedGitHubRef": stableRef,
				},
				Message:      message,
				CanRemediate: true,
			}

			issues = append(issues, &i)
		}
	}

	return issues, nil
}

func mustGetIssueMessage(workflowName string, jobName string, unstableReason UnstableReason, step workflow.ParsedWorklowStep) string {
	switch unstableReason {
	case IsStable:
		return ""
	case UnknownReason:
		return "unknown reason"
	case UnsupportedRef:
		return "unsupported ref"
	case NoSpecifiedVersion:
		return "no specified version"
	case IsMaster:
		return fmt.Sprintf("The job named %q in the %q workflow is referencing an action on the master branch of the %q repo. The master branch of %q is likely to change.",
			jobName, workflowName, step.Uses, step.Uses)
	case IsBranch:
		branch := strings.Split(step.Uses, "@")[1]
		return fmt.Sprintf("The job named %q in the %q workflow is using an action from %q. This is unstable because %q is a branch, and the contents might change.",
			jobName, workflowName, step.Uses, branch)
	case HasUnstableHistory:
		return "has unsatable history"
	case TagNotFound:
		return "tag not found"
	}

	return ""
}

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
		if path == "" {
			updatedRef = fmt.Sprintf("%s/%s@%s", owner, repo, possiblyStableTag.CommitSHA)
		} else {
			updatedRef = fmt.Sprintf("%s/%s/%s@%s", owner, repo, path, possiblyStableTag.CommitSHA)
		}
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

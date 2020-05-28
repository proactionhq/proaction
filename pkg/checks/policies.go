package checks

const (
	Outdated_Policy = `package proaction

buildOutput(repo, reason, recommendation) = built {
  built := {
    "owner": repo.owner,
    "repo": repo.repo,
    "reason": reason,
    "remediatedRef": recommendation.ref,
    "remediatedRefType": recommendation.refType,
    "checkType": "unstableGitHubRef",
    "workflow": repo.workflow,
    "originalLineNumber": to_number(repo.lineNumber),
    "originalContent": repo.lineContent,
    "remediatedContent": sprintf("%s/%s%s@%s", [repo.owner, repo.repo, repo.path, recommendation.ref])
  }
}

# utility function to check if the ref is in the tags.name
isRefInTags(repo) {
  tag := repo.tags[_]
  tag.name == repo.ref
}

# given a set of tags, get the commit sha of the tag name
findTagSHA(tagName, tags) = output {
  tag := tags[_]
  tag.name == tagName
  output := tag.head
}

# if already on a tag, recommend the latest tag
recommendations[output] {
  repo := input.repos[_]
  repo.refType == "tag"

  latestTag := repo.tags[0]
  currentTagSHA := findTagSHA(repo.ref, repo.tags)

  latestTag.head != currentTagSHA

  recommendation := {
    "ref": latestTag.name,
    "refType": "tag"
  }
  output := buildOutput(repo, "isOutdatedTag", recommendation)
}

# if it's a commit ref, make sure it's the latest
recommendations[output] {
  repo := input.repos[_]
  repo.refType == "commit"

  repo.ref != repo.commits[0]

  recommendation := {
    "ref": substring(repo.commits[0], 0, 7),
    "refType": "commit"
  }
  output := buildOutput(repo, "isOutdatedCommit", recommendation)
}
`
	Recommendations_Policy = `package proaction

buildOutput(repo, reason, recommendation) = built {
  built := {
    "owner": repo.owner,
    "repo": repo.repo,
    "reason": reason,
    "remediatedRef": recommendation.ref,
    "remediatedRefType": recommendation.refType,
    "checkType": "unstableGitHubRef",
    "workflow": repo.workflow,
    "originalLineNumber": to_number(repo.lineNumber),
    "originalContent": repo.lineContent,
    "remediatedContent": sprintf("%s/%s%s@%s", [repo.owner, repo.repo, repo.path, recommendation.ref])
  }
}


contains(refs, ref) {
  maybeRef := refs[_]
  maybeRef == ref
}

containsRef(desiredRefType, desiredRefs, actualRefType, actualRef) {
  desiredRefType == actualRefType
  contains(desiredRefs, actualRef)
}

# static recommendation
recommendations[output] {
  repo := input.repos[_]
  recommendations := [r | input.recommendations[i].owner == repo.owner; r := input.recommendations[i]]
  recommendations = [r | input.recommendations[i].repo == repo.repo; r := input.recommendations[i]]
  rec := recommendations[_]

  not containsRef(rec.refType, rec.refs, repo.refType, repo.ref)

  recommendation := {
    "ref": rec.refs[0],
    "refType": rec.refType
  }
  output = buildOutput(repo, "isStaticRecommendation", recommendation)
}

`
	Recommendations_test_Policy = `package proaction

test_containsRef {
  containsRef("a", ["b"], "a", "b")
}

test_containsRef {
  not containsRef("a", ["b"], "c", "b")
}

test_containsRef {
  not containsRef("a", ["1", "b"], "c", "b")
}

test_containsRef {
  containsRef("a1", ["b"], "a1", "b")
}

test_containsRef {
  not containsRef("a1", ["b"], "a1", "c")
}

test_containsRef {
  not containsRef("a1", ["b"], "c1", "d")
}
`
	Unfork_Policy = `package proaction

buildOutput(repo, reason, recommendation) = built {
  built := {
    "owner": repo.owner,
    "repo": repo.repo,
    "reason": reason,
    "remediatedRef": recommendation.ref,
    "remediatedRefType": recommendation.refType,
    "checkType": "unstableGitHubRef",
    "workflow": repo.workflow,
    "originalLineNumber": to_number(repo.lineNumber),
    "originalContent": repo.lineContent,
    "remediatedContent": sprintf("%s/%s%s@%s", [repo.owner, repo.repo, repo.path, recommendation.ref])
  }
}

# utility function to check if the ref is in the tags.name
isRefInTags(repo) {
  tag := repo.tags[_]
  tag.name == repo.ref
}

recommendLatestCommit(repo) = output {
  commit := repo.commits[0]
  output := {
    "ref": substring(commit, 0, 7),
    "refType": "commit"
  }
}

recommendLatestTag(repo) = output {
  tag := repo.tags[0]
  output := {
    "ref": tag.name,
    "refType": "tag"
  }
}

#  if forked, check if the commit we are on is merged into the parent
recommendations[output] {
  repo := input.repos[_]
  repo.refType == "tag"

  ## tags are sorted in the input
  repo.ref != repo.tags[0].name
  recommendation := recommendLatestTag(repo)
  output := buildOutput(repo, "isUnfork", recommendation)
}
`
	UnstableGithubRef_Policy = `package proaction

buildOutput(repo, reason, recommendation) = built {
  built := {
    "owner": repo.owner,
    "repo": repo.repo,
    "reason": reason,
    "remediatedRef": recommendation.ref,
    "remediatedRefType": recommendation.refType,
    "checkType": "unstableGitHubRef",
    "workflow": repo.workflow,
    "originalLineNumber": to_number(repo.lineNumber),
    "originalContent": repo.lineContent,
    "remediatedContent": sprintf("%s/%s%s@%s", [repo.owner, repo.repo, repo.path, recommendation.ref])
  }
}

# utility function to check if the ref is in the tags.name
isRefInTags(repo) {
  tag := repo.tags[_]
  tag.name == repo.ref
}

recommendLatestCommit(repo) = output {
  commit := repo.commits[0]
  output := {
    "ref": substring(commit, 0, 7),
    "refType": "commit"
  }
}

recommendLatestTag(repo) = output {
  tag := repo.tags[0]
  output := {
    "ref": tag.name,
    "refType": "tag"
  }
}

# default branch is unstable when it's not recommended
recommendations[output] {
  repo := input.repos[_]
  repo.refType == "branch"
  repo.ref == repo.defaultBranch

  recommendation := recommendLatestCommit(repo)
  output := buildOutput(repo, "isDefaultBranch", recommendation)
}

# unstable if it's any branch, unless a recommendation is a branch
## and there are no tags
recommendations[output] {
  repo := input.repos[_]
  repo.refType == "branch"

  ## return the newest commit on the branch
  ## but this is broken and flawed because
  ## the next time this runs, it will not
  ## know about the branch, and will update
  ## to the latest commit
  ## TODO ^^

  count(repo.tags) == 0
  recommendation := recommendLatestCommit(repo)
  output := buildOutput(repo, "isBranch", recommendation)
}

# unstable if it's any branch, unless a recommendation is a branch
## and there are tags
recommendations[output] {
  repo := input.repos[_]
  repo.refType == "branch"

  count(repo.tags) > 0
  recommendation := recommendLatestTag(repo)
  output := buildOutput(repo, "isBranch", recommendation)
}
`
)

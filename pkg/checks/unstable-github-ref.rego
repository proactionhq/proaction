package proaction

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

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

#  if forked, check if the commit we are on is merged into the parent
recommendations[output] {
  repo := input.repos[_]
  repo.refType == "tag"

  ## tags are sorted in the input
  repo.ref != repo.tags[0].name
  recommendation := recommendLatestTag(repo)
  output := buildOutput(repo, "isUnfork", recommendation)
}

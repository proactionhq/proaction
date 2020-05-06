package proaction

default allow = false

allow = true {
    count(unstable) == 0
}

checkType := "unstableGitHubRef"

# unstable if it's the default branch
unstable[output] {
  repo := input.repos[_]
  repo.refType == "branch"
  repo.ref == repo.defaultBranch

  output := {
    "owner": repo.owner,
    "repo": repo.repo,
    "reason": "isDefaultBranch",
    "checkType": checkType,
    "workflow": repo.workflow,
    "originalLineNumber": to_number(repo.lineNumber),
    "originalContent": repo.lineContent,
    "remediatedContent": "TODO"
  }
}

# unstable if it's any branch, unless a recommendation is a branch
unstable[output] {
  repo := input.repos[_]
  repo.refType == "branch"

  recommendation := bestRecommend(repo)
  recommendation.refType != "branch"

  output := {
    "owner": repo.owner,
    "repo": repo.repo,
    "reason": "isBranch",
    "checkType": checkType,
    "workflow": repo.workflow,
    "originalLineNumber": to_number(repo.lineNumber),
    "originalContent": repo.lineContent,
    "remediatedContent": recommendation
  }
}

# utility function to check if the ref is in the tags.name
isRefInTags(repo) {
  tag := repo.tags[_]
  tag.name == repo.ref
}

bestRecommend(repo) = output {
  output := staticRecommend(repo)
}

# staticRecommend returns the recommendation from the manual recommendations file
staticRecommend(repo) = output {
  recommendations := [r | input.recommendations[i].owner == repo.owner; r := input.recommendations[i]]
  recommendations = [r | input.recommendations[i].repo == repo.repo; r := input.recommendations[i]]
  recommendation := recommendations[_]

  output := {
    "ref": recommendation.refs[0],
    "refType": recommendation.refType
  }
}

# unstable if tag is not found
unstable[output] {
  repo := input.repos[_]

  repo.refType == "tag"
  not isRefInTags(repo)

  output := {
    "owner": repo.owner,
    "repo": repo.repo,
    "reason": "tagNotFound",
    "checkType": checkType,
    "workflow": repo.workflow,
    "originalLineNumber": to_number(repo.lineNumber),
    "originalContent": repo.lineContent,
    "remediatedContent": "TODO"
  }
}


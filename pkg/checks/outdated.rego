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

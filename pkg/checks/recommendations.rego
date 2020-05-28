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


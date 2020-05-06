# Policy Engine

This proposal shows how a policy engine could be used to define and execute the Proaction checks, instead of hard coded implementation.

## Goals

- Quicker execution of checks due to less runtime duplicate requests made to the GitHub API
- More testable policy execution because of easy-to-mock datasets
- Ability to quickly create new checks by writing new policy files and not code

## Non Goals

- Changes to any implementation decisions that the policies (in code) today have implemented
- Maximum efficiency and guaranteed no duplication of GitHub requests

## Background

By enabling a pipeline approach that ends with policy execution controlled by a policy engine, it would be easier to create new Proaction checks without a deep understanding of the GitHub API or Proaction codebase.
A policy engine enables a common language to parse inputs, apply rules, and produce outputs.
Using this would allow that a Proaction check could be defined in a single policy file.
This also could eventually be extended to support local (even private) policy files.

To enable this, Proaction would have to support a pipelined execution that is collect -> evaluate -> remediate.
The collect step allows the check to request specific information needed as inputs.
Collect then proceeds to collect all of this information from the APIs.
The analyze is the policy execution.
The inputs are provided to each policy and executed.
The remediate step uses the policy output to automatically remediate any possible recommendations.

## High-Level Design

To start, a policy engine needs to be selected.
Two leading choices are OpenPolicyAgent and HashiCorp Sentinel.
A detailed design must be written to allow a check to specify the inputs to pass to the collect phase.
The existing checks would then be written in the policy language.

## Detailed Design

### Collect
The collector runtime must define a maximum corpus of available data.
Each check may identify a subset of parameterized data to make available as inputs.

For the initial proposal, the following categories of data are available for any repo:

repo info: owner, name, public/private, archived, forks (maximum of 100), parent (if forked), default branch, current commit
branches: a list of all branches in the repo, each with the name and list of commit shas (maximum of 100)
tags: a list of all tags in the repo, each with the name and the commit sha it points to
commits: a list of recent commit shas in the repo (maximum of 100)

In addition, the recommendations list of Actions will be available as an input.

Checks define inputs using `yq` [path expressions](https://mikefarah.gitbook.io/yq/usage/path-expressions).
Checks can feed simple data through parsers to convert items like github refs into their parts.

An example collect implementation for a check is:

```yaml
collect:
  - name: uses
    path: jobs[*].steps[*].uses
    parser: githubref
    collectors:
      - repoInfo
      - refInfo
      - branches
      - tags
      - commits
      - recommendations
```

The output of a collect phase is stored as a YAML object that can be used as input.

```yaml
repos:
  - owner: owner1
    repo: repo1
    workflowUsed:
      - filename: workflow1.yaml
        job: jobName
        stepIndex: 0
        originalLineNumber: 14
    isPublic: false
    isArchived: false
    forks:
      - other/fork
      - other2/fork2
    parent: parent/repo
    head: 1234bac
    defaultBranch: master
    commits:
      - abcdef
      - def123
      - 123def
    branches:
      - name: branch1
        head: abcdef
    tags:
      - name: v0.1.1
        head: abcd123

```

### Evaluate

The evaluate step takes the policies in the check definition, combines them with the inputs from _all_ collectors.

```yaml
evaluate:
  - name: policy1
    rego: |
      the rego policy to evaluate
  - name: policy2
    rego: |
      package proaction.unstablegithubref

      ...
```

Each evaluation should return a JSON object:

```json
{
    "workflowFile": "workflow1.yaml",
    "job": "job1",
    "stepIndex": 0,
    "originalLineNumber": 14,
    "originalContent": "aaaaa",
    "remediatedContent": "bbbbb"
}
```

### Remediate

The remediate phase parses the output of the evaluate step and either replaces contents in the files, creates a diff, a Pull Request, or whatever workflow is desired.

This is not a check-specific implementation.

## Alternatives Considered

We could continue to build each code in code.

## Security Considerations

None identified.

# Outdated Action Example

This example workflow shows how to [outdated action](https://docs.proaction.io) check works.

To run this example with only the oudated action check enabled:

```
proaction scan --check outdated-action ./examples/outdated-action/workflow.yaml 
```

The work has several steps to illustrate how the outdated action check works.

### `actions/checkout@v1`

This is a tag, and is not considered outdated. No changes will be made to this step.

### `hashicorp/terraform-github-actions@v0.7.1`

This is a tag, and is not considered outdated. No changes will be made to this step.

### `hashicorp/terraform-github-actions@271eb39`

This is a valid commit in the repo, but not the latest. The outdated action check will recommend replacing this commit with the latest.

### `hashicorp/terraform-github-actions@49492a0`

This is the latest commit in the repo (at the time of this writing). No changes will be recommended.
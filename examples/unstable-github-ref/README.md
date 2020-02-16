# Unstable GitHub Ref Example

This example workflow shows how the [unstable github ref](https://docs.proaction.io) check works.

The workflow in this example has several steps:

### `actions/checkout@v1`
This is a common action and distributed by GitHub. This workflow references the `v1` **tag** of the `actions/checkout` repo. The source for this tag is at https://github.com/actions/checkout/tree/v1. 

Tags are not guaranteed to be immutable and should generally be avoided. But there are some commonly used actions and tags that Proaction accepts as stable. The `actions/checkout` action is considered stable -- if using a tag. This will fail if using a branch.

### `synk/actions/node@master`
This workflow references the `master` **branch** of the `synk/actions` repo.

The master branch of any repo should not be used in a workflow. A workflow that referneces the master branch is likely to be non-reproducible and may produce different output today compared to yesterday, if the action was updated. This could introduce bugs and other unexpected failures and changes to a workflow.

When it finds this workflow, Proaction will recommend changing the synk/actions step to reference a commit. This is stable and reproducible. Another Proaction check (outdated action) can be used to create suggestions for accepting the latest updates, when ready.


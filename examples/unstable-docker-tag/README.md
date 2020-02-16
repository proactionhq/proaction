# Unstable Docker Tag Example

This example workflow shows how the [unstable docker tag](https://docs.proaction.io) check works.

To run this example with only the unstable docker tag check enabled:

```
proaction scan --check unstable-docker-tag ./examples/unstable-docker-tag/workflow.yaml 
```

This example workflow uses several actions:

### `actions/checkout@v1`
This is a common action and distributed by GitHub. This workflow references the `v1` **tag** of the `actions/checkout` repo. The source for this tag is at https://github.com/actions/checkout/tree/v1. 

Tags are not guaranteed to be immutable and should generally be avoided. But there are some commonly used actions and tags that Proaction accepts as stable. The `actions/checkout` action is considered stable -- if using a tag. This will fail if using a branch.

### `docker://pactfoundation/pact-cli:latest`
This workflow sues the "latest" tag from the [pactfoundation/pact-cli](https://hub.docker.com/r/pactfoundation/pact-cli/tags) Docker image. This is subject to change, and should not be used. This check will recommend that the uses statement be changed to the latest digest (docker content addressable tag), which is immutable because it references the content.
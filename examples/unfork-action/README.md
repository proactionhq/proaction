# Unfork Action Example

This example workflows shows how the [unfork action](https://docs.proaction.io) check works.

To run this example with only the unfork-action check enabled:

```
proaction scan --check unfork-action ./examples/unfork-action/workflow.yaml
```

This example workflow uses several actions:

### `actions/checkout@v1`

This is a tagged action, and will not trigger any recommendations.

### `marccampbell/kustomize-github-action@set-image`

This action references a [fork](https://github.com/marccampbell/kustomize-github-action) of the [kustomize](https://github.com/karancode/kustomize-github-action) action. The fork was created because the action did not support a specific workflow. The fork was used, but a PR was made to the upstream action, and [was accepted and merged](https://github.com/karancode/kustomize-github-action/pull/10). At this point, the workflow should go back to the upstream action, and no longer depend on an outdated fork.


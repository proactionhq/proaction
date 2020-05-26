# Proaction GitHub Action

To enable automated, regular scans of GitHub Action workflows, it's possible to run Proaction as a GitHub Action itself. This will create automated pull requests to the repo where there are changes available.

To enable this, you'll need a [Personal Access Token](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line) with `repo` and `workflow` scope. By default, the `GITHUB_TOKEN` secret doesn't have the `workflow` scope and, as a result, cannot create pull requests to any files in the `.github/workflows` directory.

## Sample Workflow

To start, create a secret in the repo named `PROACTION_TOKEN` and provide the Personal Access Token value created above. Then, create a file named `.github/workflows/proaction.yaml` with the following content:

```yaml
on:
  schedule:
    - cron:  "0 0 * * *"

jobs:
  run-proaction:
    runs-on: ubuntu-18.04
    steps:
      - uses: actions/checkout@v2

      - uses: proactionhq/proaction/action@v0.4.2
        with:
          github-token: ${{ secrets.GITHUB_TOKEN }}

      - uses: peter-evans/create-pull-request@v2
        with:
          commit-message: "[proaction] updating workflow"
          title: Updating workflow from Proaction
          token: $${{ secrets.PROACTION_TOKEN }}
```

This workflow will run daily and create pull requests with any updated found.

### Action Inputs

All inputs are *optional*. It not set, sane defaults will be applyed.

| Name | Description | Default |
|------|-------------|---------|
| `workflow-files` | A reference to the specific workflow file(s) to scan | `.github/workflows/**` |


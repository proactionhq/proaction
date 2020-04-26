package unstablegithubref

import (
	"testing"

	"github.com/proactionhq/proaction/pkg/issue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_remediateWorkflow(t *testing.T) {
	tests := []struct {
		name                      string
		beforeContent             string
		expectedRemediatedContent string
		issue                     issue.Issue
	}{
		{
			name: "replace master with sha",
			beforeContent: `name: "Example With Unstable GitHub Refs"
on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v1

      - name: Run Snyk to check for vulnerabilities
        uses: snyk/actions/node@master
        env:
          SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }`,
			expectedRemediatedContent: `name: "Example With Unstable GitHub Refs"
on: [push]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v1

      - name: Run Snyk to check for vulnerabilities
        uses: snyk/actions/node@636d44c
        env:
          SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }`,
			issue: issue.Issue{
				CheckType:  "unstable-github-ref",
				LineNumber: 11,
				CheckData: map[string]interface{}{
					"jobName":             "build",
					"unstableReason":      IsMaster,
					"originalGitHubRef":   "snyk/actions/node@master",
					"remediatedGitHubRef": "snyk/actions/node@636d44c",
				},
				Message:      "",
				CanRemediate: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := require.New(t)

			actual, err := RemediateIssue(test.beforeContent, &test.issue)
			req.NoError(err)

			assert.Equal(t, test.expectedRemediatedContent, actual)
		})
	}
}

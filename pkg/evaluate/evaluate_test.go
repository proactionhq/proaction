package evaluate

import (
	"testing"

	"github.com/proactionhq/proaction/pkg/evaluate/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_runPolicy(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		policy      string
		input       []byte
		expectAllow bool
		expect      []types.EvaluateResult
	}{
		{
			name:  "proaction.rego",
			query: "data.proaction.recommendations",
			input: []byte(`{
    "repos": [
        {
            "id": 197814629,
            "owner": "actions",
            "repo": "checkout",
            "isPublic": true,
            "isArchived": false,
            "defaultBranch": "master",
            "isFork": false,
            "head": "",
            "ref": "master",
            "refType": "branch",
            "branches": [],
            "tags": [],
            "commits": []
        }
      ]
}

`),
			policy: `package proaction

# unstable if it's the default branch
recommendations[output] {
  repo := input.repos[_]
  repo.refType == "branch"
  repo.ref == repo.defaultBranch

  output := {
    "owner": repo.owner,
    "repo": repo.repo,
    "reason": "isDefaultBranch"
  }
}

`,
			expectAllow: false,
			expect: []types.EvaluateResult{
				{
					Owner:  "actions",
					Repo:   "checkout",
					Reason: "isDefaultBranch",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := require.New(t)

			actual, err := runPolicy(test.name, test.query, test.policy, test.input)
			req.NoError(err)

			assert.Equal(t, test.expect, actual)

			actualAllow := len(actual) == 0
			assert.Equal(t, test.expectAllow, actualAllow)
		})
	}
}

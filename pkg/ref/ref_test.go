package ref

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RefToParts(t *testing.T) {
	tests := []struct {
		ref           string
		expectedOwner string
		expectedRepo  string
		expectedPath  string
		expectedRef   string
	}{
		{
			ref:           "actions/checkout@v1",
			expectedOwner: "actions",
			expectedRepo:  "checkout",
			expectedPath:  "",
			expectedRef:   "v1",
		},
		{
			ref:           "synk/actions/node@master",
			expectedOwner: "synk",
			expectedRepo:  "actions",
			expectedPath:  "node",
			expectedRef:   "master",
		},
	}

	for _, test := range tests {
		t.Run(test.ref, func(t *testing.T) {
			req := require.New(t)

			actualOwner, actualRepo, actualPath, actualRef, err := RefToParts(test.ref)
			req.NoError(err)

			assert.Equal(t, test.expectedOwner, actualOwner)
			assert.Equal(t, test.expectedRepo, actualRepo)
			assert.Equal(t, test.expectedPath, actualPath)
			assert.Equal(t, test.expectedRef, actualRef)
		})
	}
}

func Test_DetermineGitHubRefType(t *testing.T) {
	tests := []struct {
		name                      string
		owner                     string
		repo                      string
		tag                       string
		expectedPossiblyStableTag *PossiblyStableTag
		expectedBranch            *Branch
		expectedIsCommit          bool
	}{
		{
			name:                      "a commit sha",
			owner:                     "hashicorp",
			repo:                      "terraform-github-actions",
			tag:                       "271eb39",
			expectedPossiblyStableTag: nil,
			expectedBranch:            nil,
			expectedIsCommit:          true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := require.New(t)

			actualPossiblyStableTag, actualBranch, actualIsCommit, err := DetermineGitHubRefType(test.owner, test.repo, test.tag)
			req.NoError(err)

			assert.Equal(t, test.expectedPossiblyStableTag, actualPossiblyStableTag)
			assert.Equal(t, test.expectedBranch, actualBranch)
			assert.Equal(t, test.expectedIsCommit, actualIsCommit)
		})
	}
}

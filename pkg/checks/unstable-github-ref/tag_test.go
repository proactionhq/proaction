package unstablegithubref

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getUpdatedRefFromTag(t *testing.T) {
	tests := []struct {
		name      string
		owner     string
		repo      string
		path      string
		tagName   string
		commitSHA string
		allTags   []string
		expect    string
	}{
		{
			name:      "in recommendations, is a recommended tag",
			owner:     "actions",
			repo:      "checkout",
			path:      "",
			tagName:   "v1",
			commitSHA: "123",
			allTags:   []string{},
			expect:    "actions/checkout@v1",
		},
		{
			name:      "in recommendations, is a recommended tag",
			owner:     "actions",
			repo:      "checkout",
			path:      "",
			tagName:   "v1",
			commitSHA: "123",
			allTags:   []string{},
			expect:    "actions/checkout@v1",
		},
		{
			name:      "in recommendations, not a recommended tag",
			owner:     "actions",
			repo:      "checkout",
			path:      "",
			tagName:   "1.0",
			commitSHA: "123",
			allTags: []string{
				"v1",
				"v2",
			},
			expect: "actions/checkout@v2",
		},
		{
			name:      "not in recommendations, is semver, has a specific tag available",
			owner:     "o",
			repo:      "r",
			path:      "p",
			tagName:   "v1",
			commitSHA: "123",
			allTags: []string{
				"v1",
				"v1.0.0",
				"v1.0.1",
				"v2",
			},
			expect: "o/r/p@v1.0.1",
		},
		{
			name:      "not in recommendations, not a semver will return sha",
			owner:     "o",
			repo:      "r",
			path:      "",
			tagName:   "stable",
			commitSHA: "123abc",
			allTags: []string{
				"v1",
				"v2",
			},
			expect: "o/r@123abc",
		},
		{
			name:      "not in recommendations, not a semver, has a path, will return sha",
			owner:     "o",
			repo:      "r",
			path:      "p",
			tagName:   "stable",
			commitSHA: "123abc",
			allTags: []string{
				"v1",
				"v2",
			},
			expect: "o/r/p@123abc",
		},
	}

	// set some static recommendations
	recommendationActionRefs = map[string]*Recommended{
		"actions/checkout": &Recommended{
			RecommendedRefType: "tag",
			RecommendedRefs: []string{
				"v1",
				"v2",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := require.New(t)

			actual, err := getUpdatedRefFromTag(test.owner, test.repo, test.path, test.tagName, test.commitSHA, test.allTags)
			req.NoError(err)

			assert.Equal(t, test.expect, actual)
		})
	}
}

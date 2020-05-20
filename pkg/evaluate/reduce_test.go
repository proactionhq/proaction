package evaluate

import (
	"testing"

	"github.com/proactionhq/proaction/pkg/evaluate/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ReduceResults(t *testing.T) {
	tests := []struct {
		name   string
		input  []types.EvaluateResult
		expect []types.EvaluateResult
	}{
		{
			name:   "empty",
			input:  []types.EvaluateResult{},
			expect: []types.EvaluateResult{},
		},
		{
			name: "static recommendation",
			input: []types.EvaluateResult{
				{
					Workflow:           "w",
					Owner:              "o",
					Repo:               "r",
					OriginalLineNumber: 10,
					Reason:             "blargh",
				},
				{
					Workflow:           "w",
					Owner:              "o",
					Repo:               "r",
					OriginalLineNumber: 10,
					Reason:             "isStaticRecommendation",
				},
			},
			expect: []types.EvaluateResult{
				{
					Workflow:           "w",
					Owner:              "o",
					Repo:               "r",
					OriginalLineNumber: 10,
					Reason:             "isStaticRecommendation",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := require.New(t)

			actual, err := ReduceResults(test.input)
			req.NoError(err)

			assert.Equal(t, test.expect, actual)
		})
	}
}

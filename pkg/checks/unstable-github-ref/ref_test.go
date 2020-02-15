package unstablegithubref

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_refToParts(t *testing.T) {
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

			actualOwner, actualRepo, actualPath, actualRef, err := refToParts(test.ref)
			req.NoError(err)

			assert.Equal(t, test.expectedOwner, actualOwner)
			assert.Equal(t, test.expectedRepo, actualRepo)
			assert.Equal(t, test.expectedPath, actualPath)
			assert.Equal(t, test.expectedRef, actualRef)
		})
	}
}

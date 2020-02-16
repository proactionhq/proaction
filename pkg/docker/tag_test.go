package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseImageName(t *testing.T) {
	tests := []struct {
		imageName        string
		expectedHostname string
		expectedImage    string
		expectedTag      string
	}{
		{
			imageName:        "postgres:latest",
			expectedHostname: "index.docker.io",
			expectedImage:    "library/postgres",
			expectedTag:      "latest",
		},
		{
			imageName:        "kots/kotsadm:1.11.1",
			expectedHostname: "index.docker.io",
			expectedImage:    "kots/kotsadm",
			expectedTag:      "1.11.1",
		},
		// {
		// 	imageName:           "registry.somebigbank.com:8000/myapp/myimage",
		// 	expectedHostname: "registry.somebigbank.com:8000",
		// 	expectedImage: "myapp/myimage",
		// 	expectedTag: "latest",
		// },
	}

	for _, test := range tests {
		t.Run(test.imageName, func(t *testing.T) {
			req := require.New(t)

			actualHostname, actualImage, actualTag, err := ParseImageName(test.imageName)
			req.NoError(err)

			assert.Equal(t, test.expectedHostname, actualHostname)
			assert.Equal(t, test.expectedImage, actualImage)
			assert.Equal(t, test.expectedTag, actualTag)
		})
	}
}

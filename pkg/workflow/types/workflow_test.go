package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

var (
	UbuntuLatest = "ubuntu-latest"
)

func Test_Workflow_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  GitHubWorkflow
	}{
		{
			name: "kitchen-sink",
			input: `name: "Build, test and deploy"
on: [push]

jobs:
  build-fixtures:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v1
      - uses: actions/setup-node@v1
        with:
          node-version: "10.x"
      - uses: azure/docker-login@v1
        with:
          username: proactionbot
          password: ${{ secrets.DOCKERHUB_PASSWORD }}
      - run: REPOSITORY=${{ github.repository }} GITHUB_TOKEN=${{ secrets.GITHUB_TOKEN }} make -C migrations/fixtures deps build schema-fixtures run publish

  `,
			want: GitHubWorkflow{
				Name: "Build, test and deploy",
				On: Trigger{
					Type: StrOrListType,
					StringOrListVal: &StringOrList{
						Type: List,
						ListVal: []string{
							"push",
						},
					},
				},
				Jobs: map[string]*Job{
					"build-fixtures": &Job{
						RunsOn: &StringOrList{
							Type:   String,
							StrVal: &UbuntuLatest,
						},
						Steps: []*Step{
							&Step{
								Uses: "actions/checkout@v1",
							},
							&Step{
								Uses: "actions/setup-node@v1",
								With: &With{
									Params: map[string]interface{}{
										"node-version": "10.x",
									},
								},
							},
							&Step{
								Uses: "azure/docker-login@v1",
								With: &With{
									Params: map[string]interface{}{
										"username": "proactionbot",
										"password": "${{ secrets.DOCKERHUB_PASSWORD }}",
									},
								},
							},
							&Step{
								Run: "REPOSITORY=${{ github.repository }} GITHUB_TOKEN=${{ secrets.GITHUB_TOKEN }} make -C migrations/fixtures deps build schema-fixtures run publish",
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := require.New(t)

			w := GitHubWorkflow{}
			err := yaml.Unmarshal([]byte(test.input), &w)
			req.NoError(err)

			assert.Equal(t, test.want, w)
		})
	}
}

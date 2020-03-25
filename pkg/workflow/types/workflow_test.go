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
      - name: "do something"
        uses: azure/docker-login@v1
        with:
          username: proactionbot
          password: ${{ secrets.DOCKERHUB_PASSWORD }}
      - run: REPOSITORY=${{ github.repository }} GITHUB_TOKEN=${{ secrets.GITHUB_TOKEN }} make -C migrations/fixtures deps build schema-fixtures run publish

  kots:
    runs-on: ubuntu-latest
    name: kots
    needs: [publish-3p, publish-api, publish-web]
    steps:
      - name: Checkout
        uses: actions/checkout@v1
  
      - name: "kustomize build api kots"
        uses: marccampbell/kustomize-github-action@set-image
        with:
          kustomize_version: "2.0.3"
          kustomize_build_dir: "api/kustomize/overlays/kots"
          kustomize_output_file: "kots/api.yaml"
          kustomize_set_image: "proaction-api=proactionhq/api:${{ github.sha }}"
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
								Name: "do something",
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
					"kots": &Job{
						Name: "kots",
						RunsOn: &StringOrList{
							Type:   String,
							StrVal: &UbuntuLatest,
						},
						Needs: &StringOrList{
							Type: List,
							ListVal: []string{
								"publish-3p",
								"publish-api",
								"publish-web",
							},
						},
						Steps: []*Step{
							&Step{
								Name: "Checkout",
								Uses: "actions/checkout@v1",
							},
							&Step{
								Name: "kustomize build api kots",
								Uses: "marccampbell/kustomize-github-action@set-image",
								With: &With{
									Params: map[string]interface{}{
										"kustomize_version":     "2.0.3",
										"kustomize_build_dir":   "api/kustomize/overlays/kots",
										"kustomize_output_file": "kots/api.yaml",
										"kustomize_set_image":   "proaction-api=proactionhq/api:${{ github.sha }}",
									},
								},
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

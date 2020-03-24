package types

type GitHubWorkflow struct {
	Name string                 `yaml:"name,omitempty"`
	On   Trigger                `yaml:"on"`
	Env  map[string]interface{} `yaml:"env,omitempty"`
	Jobs map[string]*Job        `yaml:"jobs"`
}

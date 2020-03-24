package types

type Job struct {
	Name   string                 `yaml:"name,omitempty"`
	Needs  *StringOrList          `yaml:"needs,omitempty"`
	RunsOn *StringOrList          `yaml:"runs-on,omitempty"`
	Env    map[string]interface{} `yaml:"env,omitempty"`
	If     string                 `yaml:"if,omitempty"`
	Steps  []*Step                `yaml:"steps,omitempty"`
}

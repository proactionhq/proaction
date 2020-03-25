package types

type Step struct {
	Name            string `yaml:"name,omitempty"`
	If              string `yaml:"if,omitempty"`
	Uses            string `yaml:"uses,omitempty"`
	Run             string `yaml:"run,omitempty"`
	Shell           string `yaml:"shell,omitempty"`
	With            *With  `yaml:"with,omitempty"`
	ContinueOnError *bool  `yaml:"continue-on-error,omitempty"`
	TimeoutMinutes  int    `yaml:"timeout-minutes,omitempty"`
}

type With struct {
	Params     map[string]interface{} `yaml:",inline,omitempty"`
	Args       string                 `yaml:"args,omitempty"`
	Entrypoint string                 `yaml:"entrypoint,omitempty"`
	Env        map[string]interface{} `yaml:"env,omitempty"`
}

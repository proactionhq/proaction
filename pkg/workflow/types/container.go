package types

type Container struct {
	Image    string                 `yaml:"image,omitempty"`
	Env      map[string]interface{} `yaml:"env,omitempty"`
	Ports    []string               `yaml:"ports,omitempty"`
	Volumes  []string               `yaml:"volumes,omitempty"`
	Options  []string               `yaml:"options,omitempty"`
	Services map[string]*Service    `yaml:"services,omitempty"`
}

type Service struct {
	Image   string   `yaml:"image,omitempty"`
	Ports   []string `yaml:"ports,omitempty"`
	Volumes []string `yaml:"volumes,omitempty"`
	Options []string `yaml:"options,omitempty"`
}

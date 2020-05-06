package types

// Collector represents a single collect definition
type Collector struct {
	Name       string   `yaml:"name"`
	Path       string   `yaml:"path"`
	Parser     string   `yaml:"parser"`
	Collectors []string `yaml:"collectors"`
}

type Output struct {
	Recommendations []*ProactionRecommendation `json:"recommendations,omitempty"`
	Repos           []*RepoOutput              `json:"repos,omitempty"`
}

type WorkflowInfo struct {
	Workflow    string `json:"workflow"`
	LineNumber  int    `json:"lineNumber"`
	LineContent string `json:"lineContent"`
}

type ProactionRecommendation struct {
	Owner   string   `json:"owner"`
	Repo    string   `json:"repo"`
	RefType string   `json:"refType"`
	Refs    []string `json:"refs"`
}

type RepoOutput struct {
	WorkflowInfo  `json:",inline"`
	ID            int64    `json:"id"`
	Owner         string   `json:"owner"`
	Repo          string   `json:"repo"`
	IsPublic      bool     `json:"isPublic"`
	IsArchived    bool     `json:"isArchived"`
	DefaultBranch string   `json:"defaultBranch"`
	IsFork        bool     `json:"isFork"`
	Forks         []string `json:"forks,omitempty"`
	Parent        string   `json:"parent,omitempty"`
	Head          string   `json:"head"`

	// included when refInfo
	Ref     string `json:"ref"`
	RefType string `json:"refType"`

	// included in branches
	Branches []BranchOutput `json:"branches,omitempty"`

	// included in tags
	Tags []TagOutput `json:"tags,omitempty"`

	// included in commits
	Commits []string `json:"commits,omitempty"`
}

type BranchOutput struct {
	Name string `json:"name"`
	Head string `json:"head"`
}

type TagOutput struct {
	Name string `json:"name"`
	Head string `json:"head"`
}

func (c Collector) Equals(other Collector) bool {
	return c.Parser == other.Parser &&
		c.Path == other.Path
}

func (c *Collector) Merge(other Collector) {
	uniqueCollectors := map[string]struct{}{}

	for _, collector := range c.Collectors {
		uniqueCollectors[collector] = struct{}{}
	}
	for _, collector := range other.Collectors {
		uniqueCollectors[collector] = struct{}{}
	}

	c.Collectors = []string{}
	for k := range uniqueCollectors {
		c.Collectors = append(c.Collectors, k)
	}
}

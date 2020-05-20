package types

type Remediation struct {
	StartLine         int
	OriginalContent   string
	RemediatedContent string

	BeforeWorkflow string
	AfterWorkflow  string
	WasRemediated  bool
}

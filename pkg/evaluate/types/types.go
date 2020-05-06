package types

type Evaluator struct {
	Name string `yaml:"name"`
	Rego string `yaml:"rego"`
}

type EvaluateResult struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`

	CheckType string `json:"checkType"`
	Reason    string `json:"reason"`

	Workflow           string `json:"workflow"`
	Message            string `json:"message"`
	OriginalLineNumber int    `json:"originalLineNumber"`
	OriginalContent    string `json:"originalContent"`

	RemediatedContent string `json:"remediatedContent"`
	RemediatedRef     string `json:"remediatedRef"`
	RemediatedRefType string `json:"remediatedRefType"`
}

type reasonPriority int

// The order of these consts defines the order priority/
// more critical on top
// there is no migration that needs to happen when reordering these
// we don't need to include everything here because the only
// sort that matters currently is static > others
const (
	isUnfork               reasonPriority = iota
	isStaticRecommendation reasonPriority = iota

	unknown reasonPriority = iota
)

func reasonToPriority(reason string) reasonPriority {
	switch reason {
	case "isStaticRecommendation":
		return isStaticRecommendation
	default:
		return unknown
	}
}

type ByPriority []EvaluateResult

func (a ByPriority) Len() int {
	return len(a)
}

func (a ByPriority) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByPriority) Less(i, j int) bool {
	return reasonToPriority(a[i].Reason) < reasonToPriority(a[j].Reason)
}

type refTypePriority int

// When there are multiple recommendations, which should we choose
// this sort defines it with the higher priority at the top
// there is no need to migrate anything if these move
const (
	refTypeTag    refTypePriority = iota
	refTypeBranch refTypePriority = iota
	refTypeCommit refTypePriority = iota

	refTypeUnknown refTypePriority = iota
)

func refTypeToPriority(refType string) refTypePriority {
	switch refType {
	case "tag":
		return refTypeTag
	case "branch":
		return refTypeBranch
	case "commit":
		return refTypeBranch
	default:
		return refTypeUnknown
	}
}

type ByRefType []EvaluateResult

func (a ByRefType) Len() int {
	return len(a)
}

func (a ByRefType) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByRefType) Less(i, j int) bool {
	return refTypeToPriority(a[i].RemediatedRefType) < refTypeToPriority(a[j].RemediatedRefType)
}

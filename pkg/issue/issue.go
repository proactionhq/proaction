package issue

import "time"

// Issue represents a single detected issue in a workflow
type Issue struct {
	ID                    string
	WorkflowID            int64
	JobName               string
	StepIdx               int
	LineNumber            int
	WorkflowMarshalledSHA []string
	CreatedAt             time.Time
	CheckType             string
	CheckData             map[string]interface{}
	Message               string
	CanRemediate          bool
	CurrentStatus         string
}

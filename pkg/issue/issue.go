package issue

import "time"

type Issue struct {
	ID                    string
	WorkflowID            int64
	JobName               string
	StepIdx               int
	WorkflowMarshalledSHA []string
	CreatedAt             time.Time
	CheckType             string
	CheckData             map[string]interface{}
	Message               string
	CanRemediate          bool
	CurrentStatus         string
}

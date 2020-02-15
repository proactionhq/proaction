package issue

import "time"

type Issue struct {
	ID                    string
	WorkflowID            int64
	WorkflowMarshalledSHA []string
	CreatedAt             time.Time
	CheckType             string
	CheckData             map[string]interface{}
	Message               string
	CanRemediate          bool
	CurrentStatus         string
}

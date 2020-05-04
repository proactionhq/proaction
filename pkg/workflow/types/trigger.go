package types

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Trigger struct {
	Type            TriggerType
	StringOrListVal *StringOrList
	MultiEventVal   *MultiEvent
}

type TriggerType int

const (
	StrOrListType TriggerType = iota
	MultiEventType
)

type MultiEvent struct {
	CheckRun                 *EventWithTypes    `yaml:"check_run,omitempty"`
	CheckSuite               *EventWithTypes    `yaml:"check_suite,omitempty"`
	Create                   *EventWithoutTypes `yaml:"create,omitempty"`
	Delete                   *EventWithoutTypes `yaml:"delete,omitempty"`
	Deployment               *EventWithoutTypes `yaml:"deployment,omitempty"`
	DeploymentStatus         *EventWithoutTypes `yaml:"deployment_status,omitempty"`
	Fork                     *EventWithoutTypes `yaml:"fork,omitempty"`
	Gollum                   *EventWithoutTypes `yaml:"gollum,omitempty"`
	IssueComment             *EventWithTypes    `yaml:"issue_comment,omitempty"`
	Issues                   *EventWithTypes    `yaml:"issues,omitempty"`
	Label                    *EventWithTypes    `yaml:"label,omitempty"`
	Milestone                *EventWithTypes    `yaml:"milestone,omitempty"`
	PageBuild                *EventWithoutTypes `yaml:"page_build,omitempty"`
	Project                  *EventWithTypes    `yaml:"project,omitempty"`
	ProjectCard              *EventWithTypes    `yaml:"project_card,omitempty"`
	ProjectColumn            *EventWithTypes    `yaml:"project_column,omitempty"`
	Public                   *EventWithoutTypes `yaml:"public,omitempty"`
	PullRequest              *PushPullEvent     `yaml:"pull_request,omitempty"`
	PullRequestReview        *EventWithTypes    `yaml:"pull_request_review,omitempty"`
	PullRequestReviewComment *EventWithTypes    `yaml:"pull_request_review_comment,omitempty"`
	Push                     *PushPullEvent     `yaml:"push,omitempty"`
	RegistryPackage          *EventWithTypes    `yaml:"registry_package,omitempty"`
	Release                  *EventWithTypes    `yaml:"release,omitempty"`
	Status                   *EventWithoutTypes `yaml:"status,omitempty"`
	Watch                    *EventWithTypes    `yaml:"watch,omitempty"`
	Schedule                 *ScheduleEvent     `yaml:"schedule,omitempty"`
	RepositoryDispatch       *EventWithoutTypes `yaml:"repository_dispatch,omitempty"`
}

type EventWithTypes struct {
	Types StringOrList `yaml:"types"`
}

type EventWithoutTypes struct {
}

type ScheduleEvent struct {
	Crons []Cron `yaml:"-"`
}

type Cron struct {
	CronField string `yaml:"cron"`
}

type PushPullEvent struct {
	Types          StringOrList `yaml:"types"`
	Branches       []string     `yaml:"branches,omitempty"`
	BranchesIgnore []string     `yaml:"branches-ignore,omitempty"`
	Tags           []string     `yaml:"tags,omitempty"`
	TagsIgnore     []string     `yaml:"tags-ignore,omitempty"`
	Paths          []string     `yaml:"paths,omitempty"`
	PathsIgnore    []string     `yaml:"paths-ignore,omitempty"`
}

func (t *Trigger) UnmarshalYAML(unmarshal func(interface{}) error) error {
	stringOrListTry := StringOrList{}
	multiEventTry := MultiEvent{}

	err := unmarshal(&stringOrListTry)
	if err == nil {
		t.Type = StrOrListType
		t.StringOrListVal = &stringOrListTry
		return nil
	}

	err = unmarshal(&multiEventTry)
	if err == nil {
		t.Type = MultiEventType
		t.MultiEventVal = &multiEventTry
		return nil
	}

	return errors.Wrapf(err, "unable to unmarshal as stringorlist or multievent")
}

func (me *MultiEvent) UnmarshalYAML(unmarshal func(interface{}) error) error {
	objects := make(map[string]interface{})
	err := unmarshal(&objects)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal multievent")
	}

	for key, rawVal := range objects {
		eventWithoutTypes := EventWithoutTypes{}

		eventWithTypes := EventWithTypes{}
		pushPullEvent := PushPullEvent{}

		if _, ok := rawVal.([]byte); ok {
			yaml.Unmarshal(rawVal.([]byte), &eventWithTypes)
			yaml.Unmarshal(rawVal.([]byte), &pushPullEvent)
		}

		switch key {
		case "check_run":
			me.CheckRun = &eventWithTypes
			break
		case "check_suite":
			me.CheckSuite = &eventWithTypes
			break
		case "create":
			me.Create = &eventWithoutTypes
			break
		case "delete":
			me.Delete = &eventWithoutTypes
			break
		case "deployment":
			me.Deployment = &eventWithoutTypes
			break
		case "deployment_status":
			me.DeploymentStatus = &eventWithoutTypes
			break
		case "fork":
			me.Fork = &eventWithoutTypes
			break
		case "gollum":
			me.Gollum = &eventWithoutTypes
			break
		case "issue_comment":
			me.IssueComment = &eventWithTypes
			break
		case "issues":
			me.Issues = &eventWithTypes
			break
		case "label":
			me.Label = &eventWithTypes
			break
		case "milestone":
			me.Milestone = &eventWithTypes
			break
		case "page_build":
			me.PageBuild = &eventWithoutTypes
			break
		case "project":
			me.Project = &eventWithTypes
			break
		case "project_card":
			me.ProjectCard = &eventWithTypes
			break
		case "project_column":
			me.ProjectColumn = &eventWithTypes
			break
		case "public":
			me.Public = &eventWithoutTypes
			break
		case "pull_request":
			me.PullRequest = &pushPullEvent
			break
		case "pull_request_review":
			me.PullRequestReview = &eventWithTypes
			break
		case "pull_request_review_comment":
			me.PullRequestReviewComment = &eventWithTypes
			break
		case "push":
			me.Push = &pushPullEvent
			break
		case "registry_package":
			me.RegistryPackage = &eventWithTypes
			break
		case "release":
			me.Release = &eventWithTypes
			break
		case "status":
			me.Status = &eventWithoutTypes
			break
		case "watch":
			me.Watch = &eventWithTypes
			break
		case "schedule":
			scheduleEvent := ScheduleEvent{}
			for _, v := range rawVal.([]interface{}) {
				vv := v.(map[string]interface{})
				if vvv, ok := vv["cron"]; ok {
					scheduleEvent.Crons = append(scheduleEvent.Crons, Cron{CronField: vvv.(string)})
				}
			}
			me.Schedule = &scheduleEvent
			break
		case "repository_dispatch":
			me.RepositoryDispatch = &eventWithoutTypes
			break
		}
	}

	return nil
}

package types

import (
	"github.com/pkg/errors"
	_ "gopkg.in/yaml.v2"
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

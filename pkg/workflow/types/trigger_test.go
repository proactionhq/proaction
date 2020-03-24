package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func Test_Trigger_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		wantType            TriggerType
		wantStringOrListVal *StringOrList
		wantMultiEventVal   *MultiEvent
	}{
		{
			name:     "string",
			input:    "foo",
			wantType: StrOrListType,
			wantStringOrListVal: &StringOrList{
				Type:    String,
				StrVal:  &foo,
				ListVal: nil,
			},
			wantMultiEventVal: nil,
		},
		{
			name:     "list",
			input:    "[foo, baz]",
			wantType: StrOrListType,
			wantStringOrListVal: &StringOrList{
				Type:   List,
				StrVal: nil,
				ListVal: []string{
					"foo",
					"baz",
				},
			},
			wantMultiEventVal: nil,
		},
		{
			name:                "create",
			input:               "create:",
			wantType:            MultiEventType,
			wantStringOrListVal: nil,
			wantMultiEventVal: &MultiEvent{
				Create: &EventWithoutTypes{},
			},
		},
		// 		{
		// 			name: "schedule",
		// 			input: `schedule:
		//   - cron: "0 */4 * * *"`,
		// 			wantType:            MultiEventType,
		// 			wantStringOrListVal: nil,
		// 			wantMultiEventVal: &MultiEvent{
		// 				Schedule: &ScheduleEvent{
		// 					Crons: []Cron{
		// 						Cron{
		// 							CronField: "0 */4 * * *",
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := require.New(t)

			tr := Trigger{}
			err := yaml.Unmarshal([]byte(test.input), &tr)
			req.NoError(err)

			assert.Equal(t, test.wantType, tr.Type)
			assert.Equal(t, test.wantStringOrListVal, tr.StringOrListVal)
			assert.Equal(t, test.wantMultiEventVal, tr.MultiEventVal)
		})
	}
}

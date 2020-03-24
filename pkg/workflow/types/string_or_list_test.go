package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

var (
	foo = "foo"
)

func Test_StringOrList_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantType    StringOrListType
		wantStrVal  *string
		wantListVal []string
	}{
		{
			name:        "string",
			input:       "foo",
			wantType:    String,
			wantStrVal:  &foo,
			wantListVal: nil,
		},
		{
			name:        "list",
			input:       `["a", "b"]`,
			wantType:    List,
			wantStrVal:  nil,
			wantListVal: []string{"a", "b"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := require.New(t)

			sl := StringOrList{}
			err := yaml.Unmarshal([]byte(test.input), &sl)
			req.NoError(err)

			assert.Equal(t, test.wantType, sl.Type)
			assert.Equal(t, test.wantStrVal, sl.StrVal)
			assert.Equal(t, test.wantListVal, sl.ListVal)
		})
	}
}

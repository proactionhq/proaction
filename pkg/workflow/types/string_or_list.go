package types

import "github.com/pkg/errors"

type StringOrList struct {
	Type    StringOrListType
	StrVal  *string
	ListVal []string
}

type StringOrListType int

const (
	String StringOrListType = iota
	List
)

func (sl *StringOrList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	strTry := ""
	listTry := []string{}

	err := unmarshal(&strTry)
	if err == nil {
		sl.Type = String
		sl.StrVal = &strTry
		return nil
	}

	err = unmarshal(&listTry)
	if err == nil {
		sl.Type = List
		sl.ListVal = *&listTry
		return nil
	}

	return errors.Wrapf(err, "unable to unmarshal as string or list")
}

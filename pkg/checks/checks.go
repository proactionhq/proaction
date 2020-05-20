package checks

import (
	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/checks/types"
)

//go:generate go run scripts/policies.go

func FromString(checkType string) (*types.Check, error) {
	switch checkType {
	case "unstable-github-ref":
		return UnstableGitHubRef(), nil
	case "outdated-action":
		return Outdated(), nil
	case "unfork-action":
		return Unfork(), nil
	default:
		return nil, errors.Errorf("unknown check type: %q", checkType)
	}
}

package types

import (
	collecttypes "github.com/proactionhq/proaction/pkg/collect/types"
	evaluatetypes "github.com/proactionhq/proaction/pkg/evaluate/types"
)

type Check struct {
	Collectors []collecttypes.Collector  `yaml:"collect"`
	Evaluators []evaluatetypes.Evaluator `yaml:"evaluate"`
}

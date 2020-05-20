package collect

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/collect/types"
)

// Collect will run the collector on the workflowContent and return the outputs
func Collect(collector types.Collector, workflowName string, workflowContent []byte) (*types.Output, error) {
	workflowInfos, err := pathsToInput(collector.Path, workflowName, workflowContent)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert paths to input")
	}

	combinedOutputs := types.Output{
		Repos: []*types.RepoOutput{},
	}
	for _, workflowInfo := range workflowInfos {
		output, err := parseInput(collector.Parser, workflowInfo, collector.Collectors)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse input")
		}

		combinedOutputs.Repos = append(combinedOutputs.Repos, output.Repos...)
	}

	// Inject the proaction collectors
	recommendations, err := loadRecommendations()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load recommendations")
	}

	combinedOutputs.Recommendations = recommendations

	return &combinedOutputs, nil
}

func parseInput(parser string, workflowInfo types.WorkflowInfo, collectors []string) (*types.Output, error) {
	if parser == "githubref" {
		output, err := parseGitHubRef(workflowInfo, collectors)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse github ref")
		}

		return output, nil
	}

	return nil, errors.New("unknown parser")
}

func loadRecommendations() ([]*types.ProactionRecommendation, error) {
	if err := loadRecommendationsFromWeb(); err != nil {
		return nil, errors.Wrap(err, "failed to load recommendations")
	}

	proactionRecommendations := []*types.ProactionRecommendation{}

	for repo, recommended := range recommendationActionRefs {
		ownerAndRepo := strings.Split(repo, "/")

		proactionRecommendation := types.ProactionRecommendation{
			Owner:   ownerAndRepo[0],
			Repo:    ownerAndRepo[1],
			RefType: recommended.RecommendedRefType,
			Refs:    recommended.RecommendedRefs,
		}

		proactionRecommendations = append(proactionRecommendations, &proactionRecommendation)
	}

	return proactionRecommendations, nil
}

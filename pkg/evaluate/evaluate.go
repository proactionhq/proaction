package evaluate

import (
	"bytes"

	"github.com/docker/distribution/context"
	"github.com/docker/go/canonical/json"
	"github.com/mitchellh/mapstructure"
	"github.com/open-policy-agent/opa/rego"
	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/evaluate/types"
)

// Evaluate will take the evaluator and input and return the standard output
func Evaluate(evaluator types.Evaluator, input []byte) ([]types.EvaluateResult, error) {
	values, err := runPolicy("proaction", "data.proaction.recommendations", evaluator.Rego, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to run policy")
	}

	return values, nil
}

func runPolicy(name string, query string, policy string, input []byte) ([]types.EvaluateResult, error) {
	r := rego.New(
		rego.Query(query),
		rego.Module(name, policy),
	)

	q, err := r.PrepareForEval(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare for eval")
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare query for eval")
	}

	d := json.NewDecoder(bytes.NewBuffer(input))
	d.UseNumber()
	var i interface{}
	if err := d.Decode(&i); err != nil {
		return nil, errors.Wrap(err, "failed to decode")
	}

	results, err := q.Eval(context.Background(), rego.EvalInput(i))
	if err != nil {
		return nil, errors.Wrap(err, "failed to eval policy")
	}

	if len(results) == 0 {
		return nil, nil
	}

	result := results[0]
	if len(result.Expressions) == 0 {
		return nil, nil
	}

	evaluateResults := []types.EvaluateResult{}

	if err := mapstructure.Decode(result.Expressions[0].Value, &evaluateResults); err != nil {
		return nil, errors.Wrap(err, "failed to mapstructure to evaluate results")
	}

	// TODO Add messages

	return evaluateResults, nil
}

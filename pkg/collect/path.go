package collect

import (
	"bytes"
	"io"

	"github.com/mikefarah/yq/v3/pkg/yqlib"
	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/collect/types"
	yaml "gopkg.in/yaml.v3"
)

var (
	yq = yqlib.NewYqLib()
)

func pathsToInput(path string, workflowName string, workflowContent []byte) ([]types.WorkflowInfo, error) {
	matchingNodes, err := findMatchingNodes(workflowContent, path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find matching nodes")
	}

	workflowInfos := []types.WorkflowInfo{}
	for _, matchingNode := range matchingNodes {
		workflowInfo := types.WorkflowInfo{
			Workflow:    workflowName,
			LineNumber:  matchingNode.Node.Line,
			LineContent: matchingNode.Node.Value,
		}

		workflowInfos = append(workflowInfos, workflowInfo)
	}

	return workflowInfos, nil
}

func findMatchingNodes(data []byte, path string) ([]*yqlib.NodeContext, error) {
	readFn := createReadFunction(path)
	matchingNodes := []*yqlib.NodeContext{}

	currentIndex := 0
	errorReadingStream := readStream(data, func(decoder *yaml.Decoder) error {
		for {
			dataBucket := yaml.Node{}
			errorReading := decoder.Decode(&dataBucket)

			if errorReading == io.EOF {
				return nil
			} else if errorReading != nil {
				return errorReading
			}

			var errorParsing error
			matchingNodes, errorParsing = appendDocument(matchingNodes, dataBucket, readFn, true, 0, currentIndex)
			if errorParsing != nil {
				return errorParsing
			}

			currentIndex = currentIndex + 1
		}
	})
	return matchingNodes, errorReadingStream
}

type readDataFn func(dataBucket *yaml.Node) ([]*yqlib.NodeContext, error)

func createReadFunction(path string) func(*yaml.Node) ([]*yqlib.NodeContext, error) {
	return func(dataBucket *yaml.Node) ([]*yqlib.NodeContext, error) {
		return yq.Get(dataBucket, path, true)
	}
}
func readStream(data []byte, yamlDecoder yamlDecoderFn) error {
	stream := bytes.NewReader(data)
	return yamlDecoder(yaml.NewDecoder(stream))
}

type yamlDecoderFn func(*yaml.Decoder) error

func appendDocument(originalMatchingNodes []*yqlib.NodeContext, dataBucket yaml.Node, readFn readDataFn, updateAll bool, docIndexInt int, currentIndex int) ([]*yqlib.NodeContext, error) {
	yqlib.DebugNode(&dataBucket)
	if !updateAll && currentIndex != docIndexInt {
		return originalMatchingNodes, nil
	}
	matchingNodes, errorParsing := readFn(&dataBucket)
	if errorParsing != nil {
		return nil, errors.Wrapf(errorParsing, "Error reading path in document index %v", currentIndex)
	}
	return append(originalMatchingNodes, matchingNodes...), nil
}

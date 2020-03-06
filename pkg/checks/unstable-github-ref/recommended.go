package unstablegithubref

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	semver "github.com/Masterminds/semver/v3"
	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/pkg/version"
)

var (
	recommendationActionRefs map[string]*Recommended
)

type Recommended struct {
	RecommendedRefType string   `json:"recommendedRefType"`
	RecommendedRefs    []string `json:"recommendedRefs"`
}

func loadRecommendationsFromWeb() error {
	masterURI := "https://raw.githubusercontent.com/proactionhq/proaction/master/pkg/checks/unstable-github-ref/recommended-action-refs.json"

	uri := ""

	currentSemver, err := semver.NewVersion(version.Version())
	if err == nil {
		uri = fmt.Sprintf("https://raw.githubusercontent.com/proactionhq/proaction/v%s/pkg/checks/unstable-github-ref/recommended-action-refs.json", currentSemver.String())
	}

	if uri != "" {
		if err := loadRecommendationsFromURI(uri); err == nil {
			if recommendationActionRefs != nil {
				return nil
			}
		}
	}

	if err := loadRecommendationsFromURI(masterURI); err != nil {
		return errors.Wrap(err, "failed to load recommendations")
	}

	return nil
}

func loadRecommendationsFromURI(uri string) error {
	res, err := http.Get(uri)
	if err != nil {
		return errors.Wrap(err, "failed to get http")
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read body")
	}

	r := map[string]*Recommended{}
	if err := json.Unmarshal(b, &r); err != nil {
		return errors.Wrap(err, "failed to parse body")
	}

	recommendationActionRefs = r
	return nil
}

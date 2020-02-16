package ref

import (
	"os"
	"strings"

	"github.com/pkg/errors"
)

// RefToParts takes a uses reference and splits into owner, repo, path and ref
func RefToParts(ref string) (string, string, string, string, error) {
	splitRef := strings.Split(ref, "@")

	if len(splitRef) < 2 {
		return "", "", "", "", errors.New("unsupported reference format")
	}

	repoParts := splitRef[0]
	tag := splitRef[1]

	splitRepoParts := strings.Split(repoParts, "/")
	owner := ""
	repo := ""
	path := ""

	if len(splitRepoParts) > 2 {
		owner = splitRepoParts[0]
		repo = splitRepoParts[1]
		path = strings.Join(splitRepoParts[2:], string(os.PathSeparator))
	} else if len(splitRepoParts) == 2 {
		owner = splitRepoParts[0]
		repo = splitRepoParts[1]
	}

	return owner, repo, path, tag, nil
}

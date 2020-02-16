package docker

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	dockerImageNameRegex = regexp.MustCompile("(?:([^\\/]+)\\/)?(?:([^\\/]+)\\/)?([^@:\\/]+)(?:[@:](.+))")
)

// ParseImageName returns the hostname, namespace, image and tag from
// a given docker image name. this doesn't validate that the image exists,
// it's only string comparison
func ParseImageName(imageName string) (string, string, string, error) {
	matches := dockerImageNameRegex.FindStringSubmatch(imageName)

	if len(matches) != 5 {
		return "", "", "", fmt.Errorf("Expected 5 matches in regex, but found %d", len(matches))
	}

	hostname := matches[1]
	namespace := matches[2]
	image := matches[3]
	tag := matches[4]

	if namespace == "" && hostname != "" {
		if !strings.Contains(hostname, ".") && !strings.Contains(hostname, ":") {
			namespace = hostname
			hostname = ""
		}
	}

	if hostname == "" {
		hostname = "index.docker.io"
	}

	if namespace == "" {
		namespace = "library"
	}

	if tag == "" {
		tag = "latest"
	}

	return hostname, fmt.Sprintf("%s/%s", namespace, image), tag, nil
}

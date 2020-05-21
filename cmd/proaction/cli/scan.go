package cli

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"

	"github.com/google/go-github/v28/github"
	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/internal/event"
	"github.com/proactionhq/proaction/pkg/checks"
	checktypes "github.com/proactionhq/proaction/pkg/checks/types"
	"github.com/proactionhq/proaction/pkg/githubapi"
	"github.com/proactionhq/proaction/pkg/logger"
	"github.com/proactionhq/proaction/pkg/scanner"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	githubPathRegex = regexp.MustCompile("/([^/?=]+)/([^/?=]+)/blob/([^/?=]+)/(.*)")
)

func ScanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "scan",
		Short:         "Check for issues in the given workflow",
		Long:          ``,
		SilenceUsage:  true,
		SilenceErrors: false,
		Args:          cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			v := viper.GetViper()

			if v.GetBool("verbose") {
				logger.SetDebug()
			}

			if err := event.Init(v); err != nil {
				if v.GetBool("debug") {
					fmt.Printf("%s\n", err.Error())
				}
			}

			files, uris, err := parseArgs(args)
			if err != nil {
				return errors.Wrap(err, "failed to parse args")
			}

			// if there are URIs, we force the "diff" flag
			// because otherwise we'd be changing tmp files
			if len(uris) > 0 {
				v.Set("diff", true)
			}

			for _, filename := range files {
				workflowContent, err := readWorkflowContentFromFile(filename)
				if err != nil {
					return errors.Wrap(err, "failed to read workflow content")
				}

				_, err = scanWorkflow(workflowContent, filename)
				if err != nil {
					return errors.Wrap(err, "failed to scan workflow")
				}
			}

			for _, uri := range uris {
				workflowContent, filename, err := readWorkflowContentFromURI(uri)
				if err != nil {
					return errors.Wrap(err, "failed to read workflow content")
				}

				_, err = scanWorkflow(workflowContent, filename)
				if err != nil {
					return errors.Wrap(err, "failed to scan workflow")
				}
			}

			return nil
		},
	}

	cmd.Flags().StringSlice("check", []string{}, "check(s) to run. if empty, all checks will run")
	cmd.Flags().String("out", "", "when set, the updated workflow will be written to the file specified, instead of in place")
	cmd.Flags().Bool("dry-run", false, "when set, proaction will print the output and recommended changes, but will not make changes to the file")
	cmd.Flags().Bool("quiet", false, "when set, proaction will not print explanations but will only update the workflow files with recommendations")
	cmd.Flags().Bool("debug", false, "when set, echo debug statements")
	cmd.Flags().Bool("diff", false, "when set, instead of writing the file, just show a diff")
	cmd.Flags().Bool("silent", false, "when set, the spinners will not be displayed")
	cmd.Flags().Bool("verbose", false, "when true, verbose logging")

	return cmd
}

// scanWorkflow will scan the workflow in content, and return a bool == true if there
// are recommendations found, and an error if we couldn't scan
func scanWorkflow(workflowContent []byte, filename string) (int, error) {
	v := viper.GetViper()

	s, err := scanner.NewScanner(filename, workflowContent)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create scanner")
	}

	if len(v.GetStringSlice("check")) == 0 {
		s.EnableAllChecks()
	} else {
		enabledChecks := []*checktypes.Check{}
		for _, requestedCheck := range v.GetStringSlice("check") {
			c, err := checks.FromString(requestedCheck)
			if err != nil {
				return 0, errors.Wrap(err, "failed to get check")
			}

			enabledChecks = append(enabledChecks, c)
		}
		s.EnableChecks(enabledChecks)
	}

	err = s.ScanWorkflow()
	if err != nil {
		return 0, errors.Wrap(err, "failed to scan workflow")
	}

	if len(s.Results) == 0 {
		fmt.Printf("no changes")
		return 0, nil
	}

	if !v.GetBool("quiet") {
		fmt.Printf("%s\n", s.GetOutput())
	}

	if !bytes.Equal(s.OriginalContent, s.RemediatedContent) {
		if v.GetBool("diff") {
			dmp := diffmatchpatch.New()
			charsA, charsB, lines := dmp.DiffLinesToChars(string(s.OriginalContent), string(s.RemediatedContent))
			diffs := dmp.DiffMain(charsA, charsB, false)
			diffs = dmp.DiffCharsToLines(diffs, lines)
			fmt.Println(dmp.DiffPrettyText(dmp.DiffCleanupEfficiency(diffs)))
		} else {
			if v.GetBool("dry-run") {
				fmt.Printf("%s\n", s.RemediatedContent)
				return 0, nil
			}

			if v.GetString("out") == "" {
				err := ioutil.WriteFile(filename, []byte(s.RemediatedContent), 0755)
				if err != nil {
					return 1, errors.Wrap(err, "failed to update workflow with remediations")
				}
			} else {
				d, _ := filepath.Split(v.GetString("out"))
				if err := os.MkdirAll(d, 0755); err != nil {
					return 1, errors.Wrap(err, "failed to mkdir for out file")
				}

				err := ioutil.WriteFile(v.GetString("out"), []byte(s.RemediatedContent), 0755)
				if err != nil {
					return 1, errors.Wrap(err, "failed to update workflow with remediations")
				}
			}
		}

		return 0, nil
	}

	return 0, nil
}

// parseArgs takes the input and returns an unglobbed list of files and urls
func parseArgs(args []string) ([]string, []string, error) {
	files := []string{}
	urls := []string{}

	for _, arg := range args {
		_, err := url.ParseRequestURI(arg)
		if err == nil {
			urls = append(urls, arg)
			continue
		}

		matches, err := filepath.Glob(arg)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to parse glob")
		}

		if matches != nil {
			files = append(files, matches...)
		}
	}

	return files, urls, nil
}

func readWorkflowContentFromURI(uri string) ([]byte, string, error) {
	// Let's be kind. If someone put in a github.com url, we can probably download
	// the file, flip a few flags around and print the recommendations to stdout
	parsedURL, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to parse uri")
	}

	localFile := ""
	// TODO we should support domains that aren't github.com
	if parsedURL.Hostname() == "github.com" {
		downloadedFile, err := downloadFileFromGitHub(parsedURL.Path)
		if err != nil {
			return nil, "", errors.Wrap(err, "tried unsuccesfully to download file from github")
		}
		defer os.RemoveAll(downloadedFile)

		localFile = downloadedFile
	} else if parsedURL.Hostname() == "raw.githubusercontent.com" {
		resp, err := http.DefaultClient.Get(parsedURL.String())
		if err != nil {
			return nil, "", errors.Wrap(err, "failed to download raw github file")
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, "", errors.Wrap(err, "failed to read response body")
		}

		tmpFile, err := ioutil.TempFile("", "proaction")
		if err != nil {
			return nil, "", errors.Wrap(err, "failed to create temp file")
		}
		defer os.RemoveAll(tmpFile.Name())

		if err := ioutil.WriteFile(tmpFile.Name(), []byte(body), 0755); err != nil {
			return nil, "", errors.Wrap(err, "failed to save to temp file")
		}

		localFile = tmpFile.Name()
	}

	content, err := ioutil.ReadFile(localFile)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to read workflow")
	}

	return content, localFile, nil
}

func readWorkflowContentFromFile(localFile string) ([]byte, error) {
	content, err := ioutil.ReadFile(localFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read workflow")
	}

	return content, nil
}

func downloadFileFromGitHub(path string) (string, error) {
	matches := githubPathRegex.FindStringSubmatch(path)

	if len(matches) != 5 {
		return "", fmt.Errorf("Expected 5 matches in regex, but found %d", len(matches))
	}

	owner := matches[1]
	repo := matches[2]
	branch := matches[3]
	filename := matches[4]

	githubClient := githubapi.NewGitHubClient()
	fileContents, _, _, err := githubClient.Repositories.GetContents(
		context.Background(), owner, repo, filename,
		&github.RepositoryContentGetOptions{
			Ref: fmt.Sprintf("heads/%s", branch),
		})
	if err != nil {
		return "", errors.Wrap(err, "failed to download contents from github")
	}

	tmpFile, err := ioutil.TempFile("", "proaction")
	if err != nil {
		return "", errors.Wrap(err, "failed to create temp file")
	}

	content, err := fileContents.GetContent()
	if err != nil {
		return "", errors.Wrap(err, "failed to get contents")
	}

	if err := ioutil.WriteFile(tmpFile.Name(), []byte(content), 0755); err != nil {
		return "", errors.Wrap(err, "failed to save to temp file")
	}

	return tmpFile.Name(), nil
}

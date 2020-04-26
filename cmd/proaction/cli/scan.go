package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/google/go-github/v28/github"
	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/internal/event"
	"github.com/proactionhq/proaction/pkg/githubapi"
	progresstypes "github.com/proactionhq/proaction/pkg/progress/types"
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

			if err := event.Init(v); err != nil {
				if v.GetBool("debug") {
					fmt.Printf("%s\n", err.Error())
				}
			}

			workflowContent, filename, err := readWorkflowContent(args)
			if err != nil {
				return errors.Wrap(err, "failed to read workflow content")
			}

			s, err := scanner.NewScanner(string(workflowContent))
			if err != nil {
				return errors.Wrap(err, "failed to create scanner")
			}

			if len(v.GetStringSlice("check")) == 0 {
				s.EnableAllChecks()
			} else {
				s.EnableChecks(v.GetStringSlice("check"))
			}

			stopChan := make(chan bool)
			stoppedChan := make(chan bool)
			go func() {
				lineCount := 0
				for {
					select {
					case <-stopChan:
						stoppedChan <- true
						return
					case <-time.After(time.Millisecond * 100):
						for i := 0; i < lineCount; i++ {
							fmt.Printf("\033[A")
						}

						maxCheckNameLength := 0
						for _, checkName := range s.EnabledChecks {
							if len(checkName) > maxCheckNameLength {
								maxCheckNameLength = len(checkName)
							}
						}

						for _, checkName := range s.EnabledChecks {
							fmt.Printf("\033[2K\r%s ", checkName)

							for i := len(checkName); i < maxCheckNameLength; i++ {
								fmt.Printf(" ")
							}

							// show the status of each check
							progress, ok := s.Progress[checkName]
							if ok {
								steps, stepStatus := progress.Get()
								for _, s := range steps {
									if status, ok := stepStatus[s]; ok {
										if status == progresstypes.ScannerStatusCompleted {
											fmt.Printf(" [%s ✓] ", s)
										} else if status == progresstypes.ScannerStatusRunning {
											fmt.Printf(" [%s ⟳] ", s)
										} else if status == progresstypes.ScannerStatusPending {
											fmt.Printf(" [%s …] ", s)
										}
									}
								}
							}
							fmt.Printf("\n")
						}

						lineCount = len(s.EnabledChecks)

					}
				}
			}()

			err = s.ScanWorkflow()
			if err != nil {
				return errors.Wrap(err, "failed to scan workflow")
			}
			stopChan <- true

			select {
			case <-stoppedChan:
				break
			case <-time.After(time.Millisecond * 100):
				break
			}
			if len(s.Issues) == 0 {
				fmt.Println("No recommendations found!")
				os.Exit(0)
			}

			if !v.GetBool("quiet") {
				fmt.Printf("%s\n", s.GetOutput())
			}

			if s.OriginalContent != s.RemediatedContent {
				if v.GetBool("diff") {
					dmp := diffmatchpatch.New()
					charsA, charsB, lines := dmp.DiffLinesToChars(s.OriginalContent, s.RemediatedContent)
					diffs := dmp.DiffMain(charsA, charsB, false)
					diffs = dmp.DiffCharsToLines(diffs, lines)
					fmt.Println(dmp.DiffPrettyText(dmp.DiffCleanupEfficiency(diffs)))
				} else {
					if v.GetBool("dry-run") {
						fmt.Printf("%s\n", s.RemediatedContent)
						return nil
					}

					if v.GetString("out") == "" {
						err := ioutil.WriteFile(filename, []byte(s.RemediatedContent), 0755)
						if err != nil {
							return errors.Wrap(err, "failed to update workflow with remediations")
						}
					} else {
						d, _ := filepath.Split(v.GetString("out"))
						if err := os.MkdirAll(d, 0755); err != nil {
							return errors.Wrap(err, "failed to mkdir for out file")
						}

						err := ioutil.WriteFile(v.GetString("out"), []byte(s.RemediatedContent), 0755)
						if err != nil {
							return errors.Wrap(err, "failed to update workflow with remediations")
						}
					}
				}

				// exit with a non-zero code because there are changes
				os.Exit(2)
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

	return cmd
}

func readWorkflowContent(args []string) ([]byte, string, error) {
	v := viper.GetViper()

	localFile := args[0]

	// Let's be kind. If someone put in a github.com url, we can probably download
	// the file, flip a few flags around and print the recommendations to stdout
	parsedURL, err := url.ParseRequestURI(localFile)
	if err == nil {
		// TODO we should support domains that aren't github.com
		if parsedURL.Hostname() == "github.com" {
			downloadedFile, err := downloadFileFromGitHub(parsedURL.Path)
			if err != nil {
				return nil, "", errors.Wrap(err, "tried unsuccesfully to download file from github")
			}
			defer os.RemoveAll(downloadedFile)

			localFile = downloadedFile

			v.Set("diff", true)
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

			v.Set("diff", true)
		}
	}

	content, err := ioutil.ReadFile(localFile)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to read workflow")
	}

	return content, localFile, nil
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

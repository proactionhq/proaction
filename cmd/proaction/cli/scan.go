package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/internal/event"
	"github.com/proactionhq/proaction/pkg/scanner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ScanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "scan",
		Short:         "r",
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

			content, err := ioutil.ReadFile(args[0])
			if err != nil {
				return errors.Wrap(err, "failed to read workflow")
			}

			s := scanner.NewScanner()
			s.OriginalContent = string(content)

			if len(v.GetStringSlice("check")) == 0 {
				s.EnableAllChecks()
			} else {
				for _, check := range v.GetStringSlice("check") {
					s.EnabledChecks = append(s.EnabledChecks, check)
				}
			}

			err = s.ScanWorkflow()
			if err != nil {
				return errors.Wrap(err, "failed to scan workflow")
			}

			if len(s.Issues) == 0 {
				fmt.Println("No recommendations found!")
				os.Exit(0)
			}

			if !v.GetBool("quiet") {
				fmt.Printf("%#v", s.GetOutput())
			}

			if s.OriginalContent != s.RemediatedContent {
				if v.GetString("out") == "" {
					err := ioutil.WriteFile(args[0], []byte(s.RemediatedContent), 0755)
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

			os.Exit(1)
			return nil
		},
	}

	cmd.Flags().StringSlice("check", []string{}, "check(s) to run. if empty, all checks will run")
	cmd.Flags().String("out", "", "when set, the updated workflow will be written to the file specified, instead of in place")
	cmd.Flags().Bool("dry-run", false, "when set, proaction will print the output and recommended changes, but will not make changes to the file")
	cmd.Flags().Bool("quiet", false, "when set, proaction will not print explanations but will only update the workflow files with recommendations")
	cmd.Flags().Bool("debug", false, "when set, echo debug statements")
	return cmd
}

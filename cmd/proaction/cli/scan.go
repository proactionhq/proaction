package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/proactionhq/proaction/internal/event"
	"github.com/proactionhq/proaction/pkg/scanner"
	"github.com/proactionhq/proaction/pkg/workflow"
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

			workflow, err := workflow.Parse(content)
			if err != nil {
				return errors.Wrap(err, "failed to parse workflow")
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

			err = s.ScanWorkflow(workflow)
			if err != nil {
				return errors.Wrap(err, "failed to scan workflow")
			}

			if len(s.Issues) == 0 {
				fmt.Println("No recommendations found!")
				os.Exit(0)
			}

			fmt.Printf("%#v", s.GetOutput())

			if s.OriginalContent != s.RemediatedContent {
				err := ioutil.WriteFile(args[0], []byte(s.RemediatedContent), 0755)
				if err != nil {
					return errors.Wrap(err, "failed to update workflow with remediations")
				}

			}

			os.Exit(1)
			return nil
		},
	}

	cmd.Flags().StringSlice("check", []string{}, "check(s) to run. if empty, all checks will run")
	cmd.Flags().Bool("dry-run", false, "when set, proaction will print the output and recommended changes, but will not make changes to the file")
	cmd.Flags().Bool("quiet", false, "when set, proaction will not print explanations but will only update the workflow files with recommendations")
	cmd.Flags().Bool("debug", false, "when set, echo debug statements")
	return cmd
}

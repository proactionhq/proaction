package cli

import (
	"fmt"

	"github.com/proactionhq/proaction/pkg/version"
	"github.com/spf13/cobra"
)

func VersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the current version and exit",
		Long:  `Print the current version and exit`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Proaction %s\n", version.Version())

			isLatest, latestVer, err := version.IsLatestRelease()
			if err != nil {
				fmt.Printf("\nUnable to check for newer releases: %s\n", err.Error())
			} else if !isLatest {
				fmt.Printf("\nVersion %s of Proaction is available.\n\nTo install the latest version, please visit https://proaction.io/docs/getting-started/installing/\n\n", latestVer)
			}

			return nil
		},
	}
	return cmd
}

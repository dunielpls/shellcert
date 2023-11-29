package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use: "version",
		Short: "Show version and exit",
		Hidden: false,
		SilenceUsage: false,
		RunE: func(cmd *cobra.command, largs []string) error {
			fmt.Println("version: todo")
			return nil
		}
	}

	rootCmd.AddCommand(cmd)
}


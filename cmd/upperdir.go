package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upperdirCmd)
}

var upperdirCmd = &cobra.Command{
	Use:   "upperdir",
	Short: "Upperdir",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

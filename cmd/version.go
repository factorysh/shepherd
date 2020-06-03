package cmd

import (
	"fmt"

	"github.com/factorysh/shepherd/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of shepherd",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Version())
	},
}

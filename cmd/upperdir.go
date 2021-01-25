package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upperdirCmd)
}

var upperdirCmd = &cobra.Command{
	Use:   "upperdir",
	Short: "Upperdir",
	RunE: func(cmd *cobra.Command, args []string) error {
		// docker
		c, err := client.NewEnvClient()
		if err != nil {
			return err
		}
		containers, err := c.ContainerList(context.TODO(), types.ContainerListOptions{})
		if err != nil {
			return err
		}
		for _, container := range containers {
			json, err := c.ContainerInspect(context.TODO(), container.ID)
			if err != nil {
				return err
			}
			u, ok := json.GraphDriver.Data["UpperDir"]
			if ok {
				fmt.Fprintln(os.Stderr, container.Names[0])
				fmt.Println(u)
			}
		}

		return nil
	},
}

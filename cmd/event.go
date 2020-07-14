package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/factorysh/shepherd/crash"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(eventCmd)
}

var eventCmd = &cobra.Command{
	Use:   "event CONTAINER",
	Short: "Displays the report for a container, as it should be sent to Sentry",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a container id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// docker
		c, err := client.NewEnvClient()
		if err != nil {
			return err
		}

		container, err := c.ContainerInspect(context.TODO(), args[0])
		if err != nil {
			return err
		}
		// crash
		cr, err := crash.New(c)
		if err != nil {
			return err
		}

		evt, err := cr.BuildEvent(&container)
		if err != nil {
			return err
		}
		j, err := json.MarshalIndent(evt, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(j))
		return nil
	},
}

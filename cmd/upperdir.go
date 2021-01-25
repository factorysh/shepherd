package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/factorysh/shepherd/du"
	"github.com/spf13/cobra"
)

var (
	all    bool
	dujson bool
)

func init() {
	rootCmd.AddCommand(upperdirCmd)
	upperdirCmd.PersistentFlags().BoolVarP(&all, "all", "a", false, "all containers")
	upperdirCmd.PersistentFlags().BoolVarP(&dujson, "json", "j", false, "json output")
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
		containers, err := c.ContainerList(context.TODO(), types.ContainerListOptions{
			All: all,
		})
		if err != nil {
			return err
		}
		if dujson {
			fmt.Println("[")
		}
		for i, container := range containers {
			json, err := c.ContainerInspect(context.TODO(), container.ID)
			if err != nil {
				return err
			}
			u, ok := json.GraphDriver.Data["UpperDir"]
			if ok {
				if dujson {
					s, err := du.Size(u)
					if err != nil {
						fmt.Fprintln(os.Stderr, err.Error())
						os.Exit(1)
					}
					fmt.Printf(`{"container":"%s", "upper_dir":"%s", "size":%d}`, container.Names[0], u, s)
					if i < len(containers)-1 {
						fmt.Println(",")
					}
				} else {
					fmt.Fprintln(os.Stderr, container.Names[0])
					fmt.Println(u)
				}
			}
		}
		if dujson {
			fmt.Println("]")
		}

		return nil
	},
}

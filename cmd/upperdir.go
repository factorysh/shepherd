package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

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

type ContainerUpperdir struct {
	Size      int64  `json:"size"`
	Nodes     int64  `json:"nodes"`
	Container string `json:"container"`
	UpperDir  string `json:"upper_dir"`
}

type ContainerUpperdirs []ContainerUpperdir

func (c ContainerUpperdirs) Len() int {
	return len(c)
}

func (c ContainerUpperdirs) Less(i, j int) bool {
	if c[i].Size != c[j].Size {
		return c[i].Size < c[j].Size
	}
	return strings.Compare(c[i].Container, c[j].Container) == -1
}

func (c ContainerUpperdirs) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

var upperdirCmd = &cobra.Command{
	Use:     "upperdir",
	Aliases: []string{"upper"},
	Short:   "List containers upperdir",
	Long: fmt.Sprintf(`
Explore content of upperdir layer of your containers.

You can pipe result : %s %s | xargs tree -s
Or using JSON output: %s %s -j | jq .`, os.Args[0], os.Args[1], os.Args[0], os.Args[1]),
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
		cs := make(ContainerUpperdirs, 0)
		for _, container := range containers {
			json, err := c.ContainerInspect(context.TODO(), container.ID)
			if err != nil {
				return err
			}
			u, ok := json.GraphDriver.Data["UpperDir"]
			if ok {
				if dujson {
					s, i, err := du.Size(u)
					if err != nil {
						fmt.Fprintln(os.Stderr, err.Error())
						os.Exit(1)
					}
					cs = append(cs, ContainerUpperdir{
						Container: container.Names[0],
						UpperDir:  u,
						Size:      s,
						Nodes:     i,
					})
				} else {
					fmt.Fprintln(os.Stderr, container.Names[0])
					fmt.Println(u)
				}
			}
		}
		if dujson {
			sort.Sort(cs)
			json.NewEncoder(os.Stdout).Encode(cs)
		}
		return nil
	},
}

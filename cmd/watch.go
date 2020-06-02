package cmd

import (
	"context"
	"time"

	"github.com/docker/docker/client"
	"github.com/factorysh/docker-visitor/visitor"
	"github.com/factorysh/janitor-go/janitor"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(watchCmd)
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch docker and clean its mess",
	Run: func(cmd *cobra.Command, args []string) {
		initLog()

		c, err := client.NewEnvClient()
		if err != nil {
			panic(err)
		}
		l := janitor.NewLater(map[string]time.Duration{"*": 10 * time.Second})
		j := janitor.New(l, c)
		w := visitor.New(c)
		w.WatchFor(j.Event)
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		err = w.Start(ctx)
		defer cancel()
		if err != nil {
			panic(err)
		}
	},
}

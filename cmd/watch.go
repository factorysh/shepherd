package cmd

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/factorysh/docker-visitor/visitor"
	"github.com/factorysh/janitor-go/config"
	"github.com/factorysh/janitor-go/janitor"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
)

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch docker and clean its mess",
	RunE: func(cmd *cobra.Command, args []string) error {
		// lorgus, sentry
		initLog()

		// docker
		c, err := client.NewEnvClient()
		if err != nil {
			return err
		}

		// config
		var cfg *config.Config

		if cfgFile != "" {
			cfg, err = config.Read(cfgFile)
			if err != nil {
				return err
			}
		} else {
			cfg = config.New()
		}

		// janitor
		l := janitor.NewLater(cfg.Ttl)
		j := janitor.New(l, c)
		w := visitor.New(c)
		w.WatchFor(j.Event)
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		err = w.Start(ctx)
		defer cancel()
		if err != nil {
			return err
		}
		return nil
	},
}

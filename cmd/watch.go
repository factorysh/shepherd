package cmd

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/factorysh/docker-visitor/visitor"
	"github.com/factorysh/shepherd/config"
	"github.com/factorysh/shepherd/crash"
	"github.com/factorysh/shepherd/metrics"
	"github.com/factorysh/shepherd/shepherd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	listenAdmin string
)

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
	watchCmd.PersistentFlags().StringVarP(&listenAdmin, "admin", "a", "localhost:4012", "Listen admin http address")
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
			log.Info(*cfg)
		} else {
			cfg = config.New()
		}

		go func() {
			err := metrics.ListenAndServe(listenAdmin)
			if err != nil {
				log.WithError(err).Error()
			}
		}()
		log.Infof("Listening http admin : http://%s/metrics", listenAdmin)
		// shepherd
		l := shepherd.NewLater(cfg.Ttl)
		j := shepherd.New(l, c)
		w := visitor.New(c)
		w.VisitCurrentCointainer(j.Visit)
		w.WatchFor(j.Event)

		// crash
		cr, err := crash.New(c)
		if err != nil {
			return err
		}
		w.WatchFor(cr.SendEvent)

		// Watch events
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

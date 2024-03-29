package cmd

import (
	"fmt"
	"os"

	"github.com/factorysh/shepherd/version"
	"github.com/onrik/logrus/filename"
	"github.com/onrik/logrus/sentry"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "shepherd",
		Short: "Shepherd cleans the mess",
		Long: fmt.Sprintf(`Shepherd watch for docker mess.
		
		version %s`, version.Version()),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func initLog() {
	filenameHook := filename.NewHook()
	log.AddHook(filenameHook)
	log.SetLevel(log.DebugLevel)
	// logrus hook for sentry, if DSN is provided
	dsn := os.Getenv("SENTRY_DSN")
	if dsn != "" {
		sentryHook, err := sentry.NewHook(sentry.Options{
			Dsn: dsn,
		}, log.PanicLevel, log.FatalLevel, log.ErrorLevel)
		if err != nil {
			panic(err)
		}
		sentryHook.AddTag("version", version.Version())
		sentryHook.AddTag("program", "shepherd")
		log.AddHook(sentryHook)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

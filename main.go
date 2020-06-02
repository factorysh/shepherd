package main

import (
	"context"
	"os"
	"time"

	"github.com/docker/docker/client"
	"github.com/factorysh/docker-visitor/visitor"
	"github.com/factorysh/janitor-go/janitor"
	"github.com/factorysh/janitor-go/version"
	"github.com/onrik/logrus/filename"
	"github.com/onrik/logrus/sentry"
	log "github.com/sirupsen/logrus"
)

func main() {
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
		sentryHook.AddTag("program", "Janitor")
		log.AddHook(sentryHook)
	}
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
}

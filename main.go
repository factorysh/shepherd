package main

import (
	"context"
	"time"

	"github.com/docker/docker/client"
	"github.com/factorysh/docker-visitor/visitor"
	"github.com/factorysh/janitor-go/janitor"
	"github.com/onrik/logrus/filename"
	log "github.com/sirupsen/logrus"
)

func main() {
	filenameHook := filename.NewHook()
	log.AddHook(filenameHook)
	log.SetLevel(log.DebugLevel)
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

package crash

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

type Crash struct {
	client *client.Client
}

func New(client *client.Client) *Crash {
	return &Crash{
		client: client,
	}
}

func (c *Crash) Event(action string, container *types.ContainerJSON) {
	fmt.Println("🦆 ", action, container.Name)
	l := log.WithField("id", container.Name)
	l.WithField("action", action).Info("crash")
	spew.Dump(container.State)
	switch action {
	case "die":
		ctx := context.TODO()
		i, err := c.client.ContainerInspect(ctx, container.ID)
		if err != nil {
			l.WithError(err).Error()
		}
		l.WithField("exitCode", i.State.ExitCode).Info()
		//spew.Dump(i)
	}
}

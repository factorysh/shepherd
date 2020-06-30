package crash

import (
	"context"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
)

type Crash struct {
	client   *client.Client
	sentry   *sentry.Client
	modifier *eventModifier
}

func New(client *client.Client) (*Crash, error) {
	c := &Crash{
		client:   client,
		modifier: &eventModifier{},
	}
	// TODO, route by project and one DSN per project or group
	dsn := os.Getenv("SENTRY_DSN")
	if dsn != "" {
		var err error
		c.sentry, err = sentry.NewClient(sentry.ClientOptions{
			// Either set your DSN here or set the SENTRY_DSN environment variable.
			Dsn: dsn,
			// Enable printing of SDK debug messages.
			// Useful when getting started or trying to figure something out.
			Debug: true,
		})
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

type eventModifier struct {
}

func (e *eventModifier) ApplyToEvent(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
	return event
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
		if c.sentry != nil {
			if i.State.ExitCode != 0 || i.State.OOMKilled {
				id := c.sentry.CaptureEvent(&sentry.Event{
					Message: fmt.Sprintf("Container crash %s", i.Config.Hostname),
					Extra: map[string]interface{}{
						"Config": i.Config,
						"State":  i.State,
					},
					Level: sentry.LevelError,
				}, nil, c.modifier)
				l.WithField("sentry", id).Info("Send to sentry")
			}
		}
	}
}

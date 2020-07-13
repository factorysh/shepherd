package crash

import (
	"context"
	"fmt"
	"io/ioutil"
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
	version  types.Version
}

func New(client *client.Client) (*Crash, error) {
	v, err := client.ServerVersion(context.TODO())
	if err != nil {
		return nil, err
	}
	c := &Crash{
		client:   client,
		modifier: &eventModifier{},
		version:  v,
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

func (c *Crash) BuildEvent(container *types.ContainerJSON) (*sentry.Event, error) {
	ctx := context.TODO()
	i, err := c.client.ContainerInspect(ctx, container.ID)
	if err != nil {
		return nil, err
	}
	ctx = context.TODO()
	r, err := c.client.ContainerLogs(ctx, container.ID, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Tail:       "250",
		Follow:     false,
	})
	defer r.Close()
	var logs []byte
	if err != nil {
		return nil, err
	}
	logs, err = ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	t := map[string]string{
		"runtime.name":    "docker",
		"runtime":         fmt.Sprintf("docker %s", c.version.Version),
		"container.name":  i.Name,
		"container.image": i.Config.Image,
	}
	v, ok := i.Config.Labels["com.docker.compose.project"]
	if ok {
		t["compose.project"] = v
	}
	v, ok = i.Config.Labels["com.docker.compose.service"]
	if ok {
		t["compose.service"] = v
	}
	return &sentry.Event{
		Message: fmt.Sprintf("Container crash %s", i.Name),
		Extra: map[string]interface{}{
			"Id":              i.ID,
			"Config":          i.Config,
			"State":           i.State,
			"GraphDriver":     i.GraphDriver,
			"Image":           i.Image,
			"Logs":            string(logs),
			"Mounts":          i.Mounts,
			"NetworkSettings": i.NetworkSettings,
			"Cgroup":          NewCgroup().fetchCgroupStates(container),
		},
		Tags:  t,
		Level: sentry.LevelError,
	}, nil
}

// SendEvent sends event to Sentry
func (c *Crash) SendEvent(action string, container *types.ContainerJSON) {
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
				evt, err := c.BuildEvent(container)
				if err != nil {
					l.WithError(err).Error()
					return
				}
				id := c.sentry.CaptureEvent(evt, nil, c.modifier)
				l.WithField("sentry", id).Info("Send to sentry")
			}
		}
	}
}

package janitor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/factorysh/janitor-go/todo"
	log "github.com/sirupsen/logrus"
)

type Janitor struct {
	later  *Later
	client *client.Client
	todo   *todo.Todo
	undead map[string]interface{}
	lock   sync.RWMutex
}

// New Janitor
func New(later *Later, client *client.Client) *Janitor {
	return &Janitor{
		later:  later,
		client: client,
		undead: make(map[string]interface{}),
		todo:   todo.New(context.Background()),
	}
}

// GetName get the name of a project
func GetName(container *types.ContainerJSON) string {
	project, ok := container.Config.Labels["com.docker.compose.project"]
	if !ok {
		return ""
	}
	return project
}

// GetTTL return the duration of an exited project
func (j *Janitor) GetTTL(name string) (time.Duration, error) {
	if name == "" {
		return j.later.Default(), nil
	}
	return j.later.Get(name)
}

// Event handle an event, from docker-visitor
func (j *Janitor) Event(action string, container *types.ContainerJSON) {
	fmt.Println("🐳 ", action)
	spew.Dump(container.State)
	l := log.WithField("id", container.ID)
	switch action {
	case "die":
		j.lock.Lock()
		j.undead[container.ID] = new(interface{})
		j.lock.Unlock()
		d, err := j.GetTTL(GetName(container))
		// Don't bother with errors, just use default duration
		if err != nil {
			l.Error(err)
			d = j.later.Default()
		}
		j.todo.Add(func() {
			err := j.client.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{})
			if err != nil {
				l.Error(err)
				return
			}
			l.Info("removed")
		}, d)
	case "destroy":
		j.lock.Lock()
		defer j.lock.Unlock()
		_, ok := j.undead[container.ID]
		if ok {
			delete(j.undead, container.ID)
		}
	}
}

func (j *Janitor) Visit(container *types.ContainerJSON) error {
	if container.State.Status == "exited" {
		d, err := j.GetTTL(GetName(container))
		l := log.WithField("id", container.ID)
		if err != nil {
			l.Error(err)
			d = j.later.Default()
		}
		f, err := time.Parse(time.RFC3339, container.State.FinishedAt)
		if err != nil {
			l.Error(err)
			// ok, it's a failure, but don't block
			return nil
		}
		age := time.Since(f)
		if age >= d {
			l.Info("Old exited container found")
			err := j.client.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{})
			if err != nil {
				l.Error(err)
			}
		} else {
			l.Info("Remove it later")
			j.todo.Add(func() {
				err := j.client.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{})
				if err != nil {
					l.Error(err)
					return
				}
				l.Info("removed")
			}, d-age)
		}
	}
	return nil
}

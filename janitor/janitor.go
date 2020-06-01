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

func New(later *Later, client *client.Client) *Janitor {
	return &Janitor{
		later:  later,
		client: client,
		undead: make(map[string]interface{}),
		todo:   todo.New(context.Background()),
	}
}

func (j *Janitor) Event(action string, container *types.ContainerJSON) {
	fmt.Println("🐳 ", action)
	spew.Dump(container.State)
	l := log.WithField("id", container.ID)
	switch action {
	case "die":
		j.lock.Lock()
		j.undead[container.ID] = new(interface{})
		j.lock.Unlock()
		var d time.Duration
		project, ok := container.Config.Labels["com.docker.compose.project"]
		if ok {
			var err error
			d, err = j.later.Get(project)
			if err != nil {
				l.Error(err)
			}
		} else {
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

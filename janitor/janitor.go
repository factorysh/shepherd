package janitor

import (
	"fmt"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/factorysh/janitor-go/todo"
)

type Janitor struct {
	later  *Later
	client *client.Client
	todo   *todo.Todo
	undead map[string]interface{}
	lock   sync.RWMutex
}

func New() *Janitor {
	return &Janitor{
		undead: make(map[string]interface{}),
	}
}

func (j *Janitor) Event(action string, container *types.ContainerJSON) {
	fmt.Println("🐳 ", action)
	spew.Dump(container.State)
	switch action {
	case "die":
		j.lock.Lock()
		j.undead[container.ID] = new(interface{})
		j.lock.Unlock()
		project, ok := container.Config.Labels["com.docker.compose.project"]
		if !ok {
			project = "no_se"
		}
		fmt.Println(project)
	/*
		d, err := j.later.Get(project)
		j.todo.Add(func() {
			j.client.ContainersPrune()
		})
	*/
	case "destroy":
		j.lock.Lock()
		defer j.lock.Unlock()
		_, ok := j.undead[container.ID]
		if ok {
			delete(j.undead, container.ID)
		}
	}
}

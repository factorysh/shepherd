package todo

import (
	"context"
	"sync"
	"time"
)

type Todo struct {
	lock    sync.RWMutex
	cpt     int64
	todo    chan int64
	actions map[int64]func()
	ctx     context.Context
}

func New(ctx context.Context) *Todo {
	t := &Todo{
		actions: make(map[int64]func()),
		todo:    make(chan int64, 100),
		ctx:     ctx,
	}
	go t.loop()
	return t
}

func (t *Todo) loop() {
	for {
		select {
		case <-t.ctx.Done():
			return
		case id := <-t.todo:
			t.lock.RLock()
			defer t.lock.RUnlock()
			a, ok := t.actions[id]
			if ok {
				a()
				delete(t.actions, id)
			} else {
				// it's a ghost
			}
		}
	}
}

// Add an action to do later
func (t *Todo) Add(f func(), later time.Duration) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.cpt++
	id := t.cpt
	t.actions[id] = f
	time.AfterFunc(later, func() {
		t.todo <- id
	})
}

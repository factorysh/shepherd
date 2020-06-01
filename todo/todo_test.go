package todo

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestTodo(t *testing.T) {
	ctx := context.TODO()
	todo := New(ctx)
	wait := sync.WaitGroup{}
	wait.Add(1)
	todo.Add(func() { wait.Done() }, 500*time.Millisecond)
	wait.Wait()
}

package store

import (
	"context"
	"strconv"
	"sync"
	"testing"
)

func TestTodoActor_ConcurrentAdd(t *testing.T) {
	actor := NewTodoActor()
	ctx := context.Background()

	const n = 100
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		i := i // capture loop var
		wg.Add(1)
		go func() {
			defer wg.Done()
			desc := "task " + strconv.Itoa(i)
			err := actor.Add(ctx, desc)
			if err != nil {
				t.Errorf("Add failed: %v", err)
			}
		}()
	}

	wg.Wait()

	todos := actor.GetAll()
	if len(todos) != n {
		t.Errorf("expected %d todos, got %d", n, len(todos))
	}
}

func TestTodoActor_ReadWriteParallel(t *testing.T) {
	actor := NewTodoActor()
	ctx := context.Background()

	t.Run("add", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 50; i++ {
			_ = actor.Add(ctx, "readwrite-"+strconv.Itoa(i))
		}
	})

	t.Run("getAll", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 50; i++ {
			_ = actor.GetAll()
		}
	})

	t.Run("edit", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 10; i++ {
			_ = actor.Edit(ctx, i, "edited-"+strconv.Itoa(i))
		}
	})

	t.Run("delete", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 10; i++ {
			_ = actor.Delete(ctx, i)
		}
	})
}

package store

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func TestStorage_SaveAndLoad(t *testing.T) {
	ctx := context.Background()

	tmp := filepath.Join(t.TempDir(), "test.json")
	storage := NewStorage[Todos](tmp)

	todos := NewTodos()
	todos.Add(ctx, "Save this")
	todos[0].CreatedAt = time.Unix(0, 0)

	err := storage.Save(ctx, todos)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	var loaded Todos
	err = storage.Load(ctx, &loaded)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	loaded[0].CreatedAt = time.Unix(0, 0)

	if len(loaded) != 1 || loaded[0].Description != "Save this" {
		t.Error("loaded data does not match saved data")
	}
}

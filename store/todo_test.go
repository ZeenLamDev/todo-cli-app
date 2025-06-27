package store

import (
	"context"
	"testing"
)

func TestAddTodo(t *testing.T) {
	ctx := context.Background()
	todos := NewTodos()
	todos.Add(ctx, "Test item")

	if len(todos) != 1 {
		t.Fatal("expected 1 todo after Add")
	}
	if todos[0].Description != "Test item" {
		t.Errorf("unexpected description: got %q", todos[0].Description)
	}
}

func TestEditTodo(t *testing.T) {
	ctx := context.Background()
	todos := NewTodos()
	todos.Add(ctx, "Old")

	err := todos.Edit(ctx, 0, "Updated")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if todos[0].Description != "Updated" {
		t.Errorf("edit failed: got %q", todos[0].Description)
	}
}

func TestDeleteTodo(t *testing.T) {
	ctx := context.Background()
	todos := NewTodos()
	todos.Add(ctx, "Task 1")
	todos.Add(ctx, "Task 2")

	err := todos.Delete(ctx, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(todos) != 1 || todos[0].Description != "Task 2" {
		t.Error("delete failed")
	}
}

func TestToggleTodo(t *testing.T) {
	ctx := context.Background()
	todos := NewTodos()
	todos.Add(ctx, "Task")

	statuses := []string{"Not started", "Started", "Completed", "Not started"}
	for i := 1; i < len(statuses); i++ {
		err := todos.Toggle(ctx, 0)
		if err != nil {
			t.Fatalf("toggle failed at step %d: %v", i, err)
		}
		if todos[0].Status != statuses[i] {
			t.Errorf("expected status %q, got %q", statuses[i], todos[0].Status)
		}
	}
}

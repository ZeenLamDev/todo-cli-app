package cli

import (
	"context"
	"testing"

	"todo/store"
)

func TestExecute_Add(t *testing.T) {
	ctx := context.Background()

	todos := store.NewTodos()
	cmd := CmdFlags{Add: "Task 1"}

	cmd.Execute(ctx, &todos)

	if len(todos) != 1 || todos[0].Description != "Task 1" {
		t.Errorf("Add command failed. Got: %+v", todos)
	}
}

func TestExecute_Toggle(t *testing.T) {
	ctx := context.Background()

	todos := store.NewTodos()
	todos.Add(ctx, "Toggle me")

	cmd := CmdFlags{Toggle: 0}
	cmd.Execute(ctx, &todos)

	if todos[0].Status != "Started" {
		t.Errorf("Expected 'Started', got: %s", todos[0].Status)
	}
}

func TestExecute_Delete(t *testing.T) {
	ctx := context.Background()
	todos := store.NewTodos()
	todos.Add(ctx, "Delete me")

	cmd := CmdFlags{Delete: 0}
	cmd.Execute(ctx, &todos)

	if len(todos) != 0 {
		t.Errorf("Delete failed. Got %d todos", len(todos))
	}
}

func TestExecute_Edit(t *testing.T) {
	ctx := context.Background()

	todos := store.NewTodos()
	todos.Add(ctx, "Initial")

	cmd := CmdFlags{Edit: "0:Updated task"}
	cmd.Execute(ctx, &todos)

	if todos[0].Description != "Updated task" {
		t.Errorf("Edit failed. Got: %s", todos[0].Description)
	}
}

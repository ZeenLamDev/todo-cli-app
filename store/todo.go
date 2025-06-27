package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
	"todo/logutil"

	"github.com/aquasecurity/table"
)

type Todo struct {
	Description string
	Status      string
	CreatedAt   time.Time
}

type Todos []Todo

func NewTodos() Todos {
	return Todos{}
}

func (todos *Todos) Add(ctx context.Context, description string) {
	logutil.Logger(ctx).Info("Adding todo", "description", description)
	todo := Todo{
		Description: description,
		Status:      "Not started",
		CreatedAt:   time.Now(),
	}
	*todos = append(*todos, todo)
}

func (todos *Todos) validateIndex(index int) error {
	if index < 0 || index >= len(*todos) {
		err := errors.New("Invalid index")
		fmt.Println(err)
		return err
	}

	return nil
}

func (todos *Todos) Delete(ctx context.Context, index int) error {
	logutil.Logger(ctx).Info("Deleting todo", "index", index)
	t := *todos

	if err := t.validateIndex(index); err != nil {
		return err
	}

	*todos = append(t[:index], t[index+1:]...)

	return nil
}

func (todos *Todos) Toggle(ctx context.Context, index int) error {
	logutil.Logger(ctx).Info("Changing todo status", "index", index)
	if err := todos.validateIndex(index); err != nil {
		return err
	}

	switch (*todos)[index].Status {
	case "Not started":
		(*todos)[index].Status = "Started"
	case "Started":
		(*todos)[index].Status = "Completed"
	case "Completed":
		(*todos)[index].Status = "Not started"
	default:
		(*todos)[index].Status = "Not started"
	}

	return nil
}

func (todos *Todos) Edit(ctx context.Context, index int, description string) error {
	logutil.Logger(ctx).Info("Editing todo", "index", index)
	t := *todos

	if err := t.validateIndex(index); err != nil {
		return err
	}

	t[index].Description = description

	return nil
}

func (todos *Todos) Print(ctx context.Context) {
	logutil.Logger(ctx).Info("Printing todos")
	table := table.New(os.Stdout)
	table.SetRowLines(false)
	table.SetHeaders("#", "Description", "Status", "Created At")

	for index, t := range *todos {
		table.AddRow(
			strconv.Itoa(index),
			t.Description,
			t.Status,
			t.CreatedAt.Format(time.RFC1123),
		)
	}

	table.Render()
}

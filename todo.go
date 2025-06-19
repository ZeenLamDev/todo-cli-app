package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aquasecurity/table"
)

type Todo struct {
	Description string
	Status      string
	CreatedAt   time.Time
}

type Todos []Todo

func (todos *Todos) add(description string) {
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

func (todos *Todos) delete(index int) error {
	t := *todos

	if err := t.validateIndex(index); err != nil {
		return err
	}

	*todos = append(t[:index], t[index+1:]...)

	return nil
}

func (todos *Todos) toggle(index int) error {
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

func (todos *Todos) edit(index int, description string) error {
	t := *todos

	if err := t.validateIndex(index); err != nil {
		return err
	}

	t[index].Description = description

	return nil
}

func (todos *Todos) print() {
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

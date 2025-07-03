package store

import (
	"context"
	"todo/logutil"
)

type AddTodo struct {
	Ctx         context.Context
	Description string
	Resp        chan error
}

type GetTodo struct {
	Index int
	Resp  chan Todo
}

type GetAllTodos struct {
	Resp chan Todos
}

type EditTodo struct {
	Ctx         context.Context
	Index       int
	Description string
	Resp        chan error
}

type DeleteTodo struct {
	Ctx   context.Context
	Index int
	Resp  chan error
}

type TodoActor struct {
	inbox chan any
	todos Todos
}

func NewTodoActor() *TodoActor {
	a := &TodoActor{
		inbox: make(chan any),
		todos: NewTodos(),
	}
	go a.run()
	return a
}

func (a *TodoActor) run() {
	for msg := range a.inbox {
		switch m := msg.(type) {
		case AddTodo:
			logutil.Logger(m.Ctx).Info("Adding todo", "desc", m.Description)
			a.todos.Add(m.Ctx, m.Description)
			m.Resp <- nil

		case GetTodo:
			if m.Index >= 0 && m.Index < len(a.todos) {
				m.Resp <- a.todos[m.Index]
			} else {
				m.Resp <- Todo{}
			}

		case GetAllTodos:
			m.Resp <- a.todos

		case EditTodo:
			logutil.Logger(m.Ctx).Info("Editing todo", "index", m.Index)
			err := a.todos.Edit(m.Ctx, m.Index, m.Description)
			m.Resp <- err

		case DeleteTodo:
			logutil.Logger(m.Ctx).Info("Deleting todo", "index", m.Index)
			err := a.todos.Delete(m.Ctx, m.Index)
			m.Resp <- err
		}
	}
}

func (a *TodoActor) Add(ctx context.Context, desc string) error {
	resp := make(chan error)
	a.inbox <- AddTodo{Ctx: ctx, Description: desc, Resp: resp}
	return <-resp
}

func (a *TodoActor) Get(index int) Todo {
	resp := make(chan Todo)
	a.inbox <- GetTodo{Index: index, Resp: resp}
	return <-resp
}

func (a *TodoActor) GetAll() Todos {
	resp := make(chan Todos)
	a.inbox <- GetAllTodos{Resp: resp}
	return <-resp
}

func (a *TodoActor) Edit(ctx context.Context, index int, desc string) error {
	resp := make(chan error)
	a.inbox <- EditTodo{Ctx: ctx, Index: index, Description: desc, Resp: resp}
	return <-resp
}

func (a *TodoActor) Delete(ctx context.Context, index int) error {
	resp := make(chan error)
	a.inbox <- DeleteTodo{Ctx: ctx, Index: index, Resp: resp}
	return <-resp
}

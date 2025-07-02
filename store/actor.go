package store

import "context"

type AddTodo struct {
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
	Index       int
	Description string
	Resp        chan error
}

type DeleteTodo struct {
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
			a.todos.Add(context.TODO(), m.Description)
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
			err := a.todos.Edit(context.TODO(), m.Index, m.Description)
			m.Resp <- err

		case DeleteTodo:
			err := a.todos.Delete(context.TODO(), m.Index)
			m.Resp <- err
		}
	}
}

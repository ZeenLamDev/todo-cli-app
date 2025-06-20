package main

import (
	"todo/cli"
	"todo/store"
)

func main() {
	todos := store.NewTodos()
	storage := store.NewStorage[store.Todos]("todos.json")
	storage.Load(&todos)
	cmdFlags := cli.NewCmdFlag()
	cmdFlags.Execute(&todos)
	storage.Save(todos)
}

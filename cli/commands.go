package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	store "todo/store"
)

type CmdFlags struct {
	Add    string
	Delete int
	Edit   string
	Toggle int
	List   bool
}

func NewCmdFlag() *CmdFlags {
	cf := CmdFlags{}

	flag.StringVar(&cf.Add, "add", "", "Add a new todo specify description")
	flag.StringVar(&cf.Edit, "edit", "", "Edit a todo by index & specify a new description. id:new_description")
	flag.IntVar(&cf.Delete, "delete", -1, "Specify a todo by index to delete")
	flag.IntVar(&cf.Toggle, "toggle", -1, "Specify a todo by index to cycle status")
	flag.BoolVar(&cf.List, "list", false, "List all todos")

	flag.Parse()

	return &cf

}

func (cf *CmdFlags) Execute(ctx context.Context, todos *store.Todos) {
	switch {
	case cf.List:
		todos.Print(ctx)
	case cf.Add != "":
		todos.Add(ctx, cf.Add)
	case cf.Edit != "":
		parts := strings.SplitN(cf.Edit, ":", 2)
		if len(parts) != 2 {
			fmt.Print("Error, invalid format for edit. Please use id:new_description")
			os.Exit(1)
		}
		index, err := strconv.Atoi(parts[0])

		if err != nil {
			fmt.Println("Error: invalid index for edit")
			os.Exit(1)
		}

		todos.Edit(ctx, index, parts[1])

	case cf.Toggle != -1:
		todos.Toggle(ctx, cf.Toggle)
		if err := todos.Toggle(ctx, cf.Delete); err != nil {
			fmt.Println("Error:", err)
		}
	case cf.Delete != -1:
		todos.Delete(ctx, cf.Delete)
		if err := todos.Delete(ctx, cf.Delete); err != nil {
			fmt.Println("Error:", err)
		}
	default:
		fmt.Println("Invalid command")
	}
}

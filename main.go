package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/reimarrosas/task-go/models"
	"github.com/urfave/cli/v2"
)

type Ops struct {
	task *models.Task
}

func (o *Ops) AddTask(ctx *cli.Context) error {
	args := ctx.Args()
	todo := strings.Join(args.Slice(), " ")

	if err := o.task.Add(todo); err != nil {
		return err
	}

    fmt.Printf("Added `%s` to your task list.\n", todo)

	return nil
}

func (o *Ops) DoTask(ctx *cli.Context) error {
	args := ctx.Args()

	if ctx.Args().Len() != 1 {
		return errors.New("invalid number of arguments")
	}

	if err := o.task.Delete(args.First()); err != nil {
		return err
	}

    fmt.Printf("Marked %s as completed.\n", args.First())

	return nil
}

func (o *Ops) ListTask(ctx *cli.Context) error {
	ts, err := o.task.List()
	if err != nil {
		return err
	}

	if len(ts) == 0 {
		fmt.Println("You have no task to complete!")
		return nil
	}

	fmt.Println("You have the following tasks:")
	for i, v := range ts {
		fmt.Printf("%d. %s\n", i+1, v.Description)
	}

	return nil
}

func main() {
	t, err := models.InitTask("task.db")
	if err != nil {
		log.Fatal(err)
	}

	ops := Ops{t}

	app := &cli.App{
		Name:      "task",
		Usage:     "A CLI for managing your TODOs",
		UsageText: "task [command]",
		Commands: []*cli.Command{
			{
				Name:   "add",
				Usage:  "Add a new task to your TODO list",
				Action: ops.AddTask,
			},
			{
				Name:   "do",
				Usage:  "Mark a task on your TODO list as complete",
				Action: ops.DoTask,
			},
			{
				Name:   "list",
				Usage:  "List all of your incomplete tasks",
				Action: ops.ListTask,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

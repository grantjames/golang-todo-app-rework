package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	todos "grantjames.github.io/m/v2"
	"grantjames.github.io/m/v2/logger"
)

var todosFile = "todos.json"

func main() {
	// Log to file so the console output is easier to read
	f, err := os.Create("app.log") // This will truncate the file on each run
	if err != nil {
		panic(err)
	}
	defer f.Close()

	desc := flag.String("add", "", "Description of todo item to add")
	updateId := flag.String("id", "", "ID of todo to update")
	updateDesc := flag.String("desc", "", "New description for update")
	updateStatus := flag.String("status", "", "New status for update")
	deleteId := flag.String("delete", "", "ID of todo to delete")
	flag.Parse()

	ctx := context.WithValue(context.Background(), logger.TraceIDKey{}, uuid.New().String())
	todos.StartStore(todosFile)

	if *desc != "" {
		slog.Info("CLI calling store.Create()")
		todos.Create(ctx, *desc)
	}

	if *updateId != "" && (*updateDesc != "" || *updateStatus != "") {
		slog.Info("CLI calling store.Update()")
		ok := todos.Update(ctx, *updateId, *updateDesc, *updateStatus)
		if ok {
			fmt.Printf("Updated todo #%s\n", *updateId)
		}
	}

	if *deleteId != "" {
		slog.Info("CLI calling store.Delete()")
		ok := todos.Delete(ctx, *deleteId)
		if ok {
			fmt.Printf("Deleted todo #%s\n", *deleteId)
		}
	}

	fmt.Println("Todo list:")
	slog.Info("CLI calling store.List()")
	todos := todos.List(ctx)
	for i, t := range todos {
		fmt.Printf("%s: %s [%s]\n", i, t.Description, t.Status)
	}

	// Wait for interrupt signal before saving and exiting
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Press Ctrl+C to exit.")
	<-sigs
}

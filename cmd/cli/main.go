package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	todos "grantjames.github.io/m/v2"
)

var todosFile = "todos.json"

// This is not needed if we are just attaching the trace_id to the default logger
//type traceIdKey struct{}

func main() {
	// Log to file so the console output is easier to read
	f, err := os.Create("app.log") // This will truncate the file on each run
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Generate a TraceID for this run using uuid
	traceID := uuid.New().String()
	// Does this need to be added to the context since we can just attach it to default logger?
	// Would need to pass context around only if the CLI was a REPL that sent multiple requests to the store
	//ctx := context.WithValue(context.Background(), traceIdKey{}, traceID)
	logger := slog.New(slog.NewTextHandler(f, nil))
	slog.SetDefault(logger.With("trace_id", traceID))

	desc := flag.String("add", "", "Description of todo item to add")
	updateId := flag.String("id", "", "ID of todo to update")
	updateDesc := flag.String("desc", "", "New description for update")
	updateStatus := flag.String("status", "", "New status for update")
	deleteId := flag.String("delete", "", "ID of todo to delete")
	flag.Parse()

	todos.StartStore(todosFile)

	if *desc != "" {
		slog.Info("CLI calling store.Create()")
		todos.Create(*desc)
	}

	if *updateId != "" && (*updateDesc != "" || *updateStatus != "") {
		slog.Info("CLI calling store.Update()")
		ok := todos.Update(*updateId, *updateDesc, *updateStatus)
		if ok {
			fmt.Printf("Updated todo #%s\n", *updateId)
		}
	}

	if *deleteId != "" {
		slog.Info("CLI calling store.Delete()")
		ok := todos.Delete(*deleteId)
		if ok {
			fmt.Printf("Deleted todo #%s\n", *deleteId)
		}
	}

	fmt.Println("Todo list:")
	slog.Info("CLI calling store.List()")
	todos := todos.List()
	for i, t := range todos {
		fmt.Printf("%s: %s [%s]\n", i, t.Description, t.Status)
	}

	// Wait for interrupt signal before saving and exiting
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Press Ctrl+C to exit.")
	<-sigs
}

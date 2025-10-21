package todos

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
)

type Todo struct {
	Description string `json:"description"`
	Status      string `json:"status"`
}

var filePath = ""
var todos = make(map[string]Todo)
var cmds = make(chan func())

func StartStore(file string) {
	// This needs looking in to - assumes the cwd is cmd/cli or cmd/api
	filePath = fmt.Sprintf("../../%s", file)
	todos = loadTodosFromFile(filePath)
	slog.Info("Starting todo store", "file_path", filePath)
	go func() {
		for f := range cmds {
			f()
		}
	}()
}

func Get(id string) (Todo, error) {
	slog.Info("Retrieving todo from store", "todo_id", id)

	r := make(chan struct {
		t  Todo
		ok bool
	}, 1)
	cmds <- func() {
		if t, ok := todos[id]; ok {
			r <- struct {
				t  Todo
				ok bool
			}{t, true}
		} else {
			r <- struct {
				t  Todo
				ok bool
			}{Todo{}, false}
		}
	}
	v := <-r
	if !v.ok {
		return Todo{}, errors.New("Todo was not found")
	}
	return v.t, nil
}

func List() map[string]Todo {
	slog.Info("Retrieving all todos from store")

	r := make(chan map[string]Todo, 1)
	cmds <- func() {
		r <- todos
	}
	return <-r
}

func Create(desc string) string {
	slog.Info("Creating new todo in store", "description", desc)
	r := make(chan string, 1)
	cmds <- func() {
		id := uuid.NewString()
		todos[id] = Todo{Description: desc, Status: "not started"}
		saveTodosToFile(filePath, todos)
		r <- id
	}
	return <-r
}

func Update(id string, desc string, status string) bool {
	r := make(chan bool, 1)
	cmds <- func() {
		if _, ok := todos[id]; ok {
			if desc == "" {
				desc = todos[id].Description
			}
			if status == "" {
				status = todos[id].Status
			}
			todos[id] = Todo{Description: desc, Status: status}
			saveTodosToFile(filePath, todos)
			r <- true
		} else {
			r <- false
		}
	}
	return <-r
}

func Delete(id string) bool {
	r := make(chan bool, 1)
	cmds <- func() {
		if _, ok := todos[id]; ok {
			delete(todos, id)
			saveTodosToFile(filePath, todos)
			r <- true
		} else {
			r <- false
		}
	}
	return <-r
}

func loadTodosFromFile(file string) map[string]Todo {
	var todos *map[string]Todo
	// This needs looking in to - assumes the cwd is cmd/cli or cmd/api
	f, err := os.Open(file)
	if err != nil {
		slog.Error("Could not load json file, will create a new one on save", "error", err.Error())
		return map[string]Todo{}
	}
	defer f.Close()
	json.NewDecoder(f).Decode(&todos)
	return *todos
}

func saveTodosToFile(file string, todos map[string]Todo) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(todos)
}

package todos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"grantjames.github.io/m/v2/logger"
)

type Todo struct {
	Description string `json:"description"`
	Status      string `json:"status"`
}

type Request struct {
	Name     string
	Id       string
	Todo     Todo
	Response chan Response
}

type Response struct {
	Id    string
	Todo  Todo
	Todos map[string]Todo
	Ok    bool
}

var filePath = ""
var todos = make(map[string]Todo)

var cmds = make(chan Request)

func StartStore(file string) {
	// This needs looking in to - assumes the cwd is cmd/cli or cmd/api
	filePath = fmt.Sprintf("../../%s", file)
	todos = loadTodosFromFile(filePath)
	slog.Info("Starting todo store", "file_path", filePath)
	go func() {
		for cmd := range cmds {
			switch cmd.Name {
			case "get":
				t, ok := todos[cmd.Id]
				cmd.Response <- Response{
					Ok:   ok,
					Todo: t,
				}
			case "list":
				cmd.Response <- Response{
					Ok:    true,
					Todos: todos,
				}
			case "create":
				id := uuid.NewString()
				todos[id] = Todo{Description: cmd.Todo.Description, Status: "not started"}
				saveTodosToFile(filePath, todos)
			case "update":
				if _, ok := todos[cmd.Id]; ok {
					if cmd.Todo.Description == "" {
						cmd.Todo.Description = todos[cmd.Id].Description
					}
					if cmd.Todo.Status == "" {
						cmd.Todo.Status = todos[cmd.Id].Status
					}
					todos[cmd.Id] = Todo{Description: cmd.Todo.Description, Status: cmd.Todo.Status}
					saveTodosToFile(filePath, todos)
					cmd.Response <- Response{
						Ok: true,
					}
				} else {
					cmd.Response <- Response{
						Ok: false,
					}
				}
			case "delete":
				if _, ok := todos[cmd.Id]; ok {
					delete(todos, cmd.Id)
					saveTodosToFile(filePath, todos)
					cmd.Response <- Response{
						Ok: true,
					}
				} else {
					cmd.Response <- Response{
						Ok: false,
					}
				}
			}

		}
	}()
}

func Get(ctx context.Context, id string) (Todo, error) {
	logger.ContextLogger(ctx).Info("Retrieving todo from store", "todo_id", id)

	r := make(chan Response, 1)
	cmds <- Request{
		Name:     "get",
		Id:       id,
		Response: r,
	}
	resp := <-r
	if !resp.Ok {
		return Todo{}, errors.New("Todo was not found")
	}
	return resp.Todo, nil

}

func List(ctx context.Context) map[string]Todo {
	slog.InfoContext(ctx, "Retrieving all todos from store")
	r := make(chan Response, 1)
	cmds <- Request{
		Name:     "list",
		Response: r,
	}
	resp := <-r
	if resp.Ok {
		return resp.Todos
	}

	return map[string]Todo{}
}

func Create(ctx context.Context, desc string) string {
	slog.InfoContext(ctx, "Creating new todo in store", "description", desc)
	r := make(chan Response, 1)
	cmds <- Request{
		Name: "create",
		Todo: Todo{
			Description: desc,
			Status:      "not started",
		},
		Response: r,
	}
	resp := <-r
	return resp.Id
}

func Update(ctx context.Context, id string, desc string, status string) bool {
	slog.InfoContext(ctx, "Updating todo with ID", "todo_id", id)
	r := make(chan Response, 1)
	cmds <- Request{
		Name: "update",
		Id:   id,
		Todo: Todo{
			Description: desc,
			Status:      status,
		},
		Response: r,
	}
	resp := <-r
	return resp.Ok
}

func Delete(ctx context.Context, id string) bool {
	slog.InfoContext(ctx, "Deleting todo", "todo_id", id)
	r := make(chan Response, 1)
	cmds <- Request{
		Name:     "delete",
		Id:       id,
		Response: r,
	}
	resp := <-r
	return resp.Ok
}

func loadTodosFromFile(file string) map[string]Todo {
	var todos *map[string]Todo
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

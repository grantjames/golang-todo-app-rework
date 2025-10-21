package todos

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"log/slog"
// 	"os"

// 	"github.com/google/uuid"
// )

// type Todo struct {
// 	Description string `json:"description"`
// 	Status      string `json:"status"`
// }

// type Store struct {
// 	filePath string
// 	todos    map[string]Todo
// 	cmds     chan func()
// }

// func NewStore(filepath string) *Store {
// 	// This needs looking in to - assumes the cwd is cmd/cli or cmd/api
// 	filepath = fmt.Sprintf("../../%s", filepath)
// 	s := &Store{
// 		filePath: filepath,
// 		cmds:     make(chan func()),
// 		todos:    loadTodosFromFile(filepath),
// 	}
// 	go func() {
// 		for f := range s.cmds {
// 			f()
// 		}
// 	}()
// 	return s
// }

// func (s *Store) Get(id string) (Todo, error) {
// 	slog.Info("Retrieving todo from store", "todo_id", id)

// 	r := make(chan struct {
// 		t  Todo
// 		ok bool
// 	}, 1)
// 	s.cmds <- func() {
// 		if t, ok := s.todos[id]; ok {
// 			r <- struct {
// 				t  Todo
// 				ok bool
// 			}{t, true}
// 		} else {
// 			r <- struct {
// 				t  Todo
// 				ok bool
// 			}{Todo{}, false}
// 		}
// 	}
// 	v := <-r
// 	if !v.ok {
// 		return Todo{}, errors.New("Todo was not found")
// 	}
// 	return v.t, nil
// }

// func (s *Store) List() map[string]Todo {
// 	slog.Info("Retrieving all todos from store")

// 	r := make(chan map[string]Todo, 1)
// 	s.cmds <- func() {
// 		r <- s.todos
// 	}
// 	return <-r
// }

// func (s *Store) Create(desc string) string {
// 	slog.Info("Creating new todo in store", "description", desc)
// 	r := make(chan string, 1)
// 	s.cmds <- func() {
// 		id := uuid.NewString()
// 		s.todos[id] = Todo{Description: desc, Status: "not started"}
// 		saveTodosToFile(s.filePath, s.todos)
// 		r <- id
// 	}
// 	return <-r
// }

// func (s *Store) Update(id string, desc string, status string) bool {
// 	r := make(chan bool, 1)
// 	s.cmds <- func() {
// 		if _, ok := s.todos[id]; ok {
// 			if desc == "" {
// 				desc = s.todos[id].Description
// 			}
// 			if status == "" {
// 				status = s.todos[id].Status
// 			}
// 			s.todos[id] = Todo{Description: desc, Status: status}
// 			saveTodosToFile(s.filePath, s.todos)
// 			r <- true
// 		} else {
// 			r <- false
// 		}
// 	}
// 	return <-r
// }

// func (s *Store) Delete(id string) bool {
// 	r := make(chan bool, 1)
// 	s.cmds <- func() {
// 		if _, ok := s.todos[id]; ok {
// 			delete(s.todos, id)
// 			saveTodosToFile(s.filePath, s.todos)
// 			r <- true
// 		} else {
// 			r <- false
// 		}
// 	}
// 	return <-r
// }

// func loadTodosFromFile(file string) map[string]Todo {
// 	var todos *map[string]Todo
// 	// This needs looking in to - assumes the cwd is cmd/cli or cmd/api
// 	f, err := os.Open(file)
// 	if err != nil {
// 		slog.Error("Could not load json file, will create a new one on save", "error", err.Error())
// 		return map[string]Todo{}
// 	}
// 	defer f.Close()
// 	json.NewDecoder(f).Decode(&todos)
// 	return *todos
// }

// func saveTodosToFile(file string, todos map[string]Todo) error {
// 	f, err := os.Create(file)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()
// 	enc := json.NewEncoder(f)
// 	enc.SetIndent("", "  ")
// 	return enc.Encode(todos)
// }

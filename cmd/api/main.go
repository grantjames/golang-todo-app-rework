package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"text/template"

	"github.com/google/uuid"
	todos "grantjames.github.io/m/v2"
	"grantjames.github.io/m/v2/logger"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/create/", handleCreate)
	mux.HandleFunc("/get/", handleGet)
	mux.HandleFunc("/update/", handleUpdate)
	mux.HandleFunc("/delete/", handleDelete)
	static := http.FileServer(http.Dir("../../web/static/about"))
	mux.Handle("/about/", http.StripPrefix("/about/", static))

	tmpl := template.Must(template.ParseFiles("../../web/templates/list.html"))
	mux.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		items := todos.List(r.Context())
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(w, items)
	})

	todos.StartStore("todos.json")

	handler := traceIDMiddleware(mux)
	slog.Info("API server listening on :5000")
	log.Fatal(http.ListenAndServe(":5000", handler))
}

func traceIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New().String()
		ctx := context.WithValue(r.Context(), logger.TraceIDKey{}, traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	logger.ContextLogger(r.Context()).Info("/create")
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var t todos.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	idx := todos.Create(r.Context(), t.Description)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(idx))
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	logger.ContextLogger(r.Context()).Info("/get")

	id := strings.TrimPrefix(r.URL.Path, "/get/")
	todo, err := todos.Get(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	logger.ContextLogger(r.Context()).Info("/update")

	var req struct {
		Id     string `json:"id"`
		Desc   string `json:"desc"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ok := todos.Update(r.Context(), req.Id, req.Desc, req.Status)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	logger.ContextLogger(r.Context()).Info("/delete")

	id := strings.TrimPrefix(r.URL.Path, "/delete/")
	ok := todos.Delete(r.Context(), id)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

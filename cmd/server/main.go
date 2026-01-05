// cmd/server/main.go (минимум для запуска)
package main

import (
	"log"
	"net/http"

	"tasks-api/internal/handlers"
	"tasks-api/internal/storage"
	"tasks-api/internal/storage/memory"
)

func main() {
	// TODO: подключите конкретную реализацию (in‑memory) интерфейса Storage
	var store storage.Storage = memory.New()

	h := handlers.New(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", h.TasksCollection) // GET, POST
	mux.HandleFunc("/tasks/", h.TaskItem)       // GET, PUT, DELETE

	log.Println("server listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

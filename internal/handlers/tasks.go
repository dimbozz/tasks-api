// internal/handlers/tasks.go (заглушки обработчиков)
package handlers

import (
    "net/http"
    "tasks-api/internal/storage"
)

type Handler struct{ Store storage.Storage }

func New(s storage.Storage) *Handler { return &Handler{Store: s} }

// /tasks (GET, POST)
func (h *Handler) TasksCollection(w http.ResponseWriter, r *http.Request) {
    // TODO: реализуйте разбор метода, JSON, коды статусов, валидацию
}

// /tasks/{id} (GET, PUT, DELETE)
func (h *Handler) TaskItem(w http.ResponseWriter, r *http.Request) {
    // TODO: извлечение id, маршрутизация по методу, ошибки
}

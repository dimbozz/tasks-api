// internal/handlers/tasks.go (заглушки обработчиков)
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"tasks-api/internal/storage"
)

type Handler struct{ Store storage.Storage }

func New(s storage.Storage) *Handler { return &Handler{Store: s} }

// /tasks (GET, POST)
func (h *Handler) TasksCollection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// Список задач
		tasks := h.Store.List()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tasks)

	case http.MethodPost:
		// Создать задачу (заглушка ШАГ 4)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":      1,
			"message": "Task created (step 4)",
		})

	default:
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// /tasks/{id} (GET, PUT, DELETE)
func (h *Handler) TaskItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Извлечение ID из пути
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, `{"error": "Task ID required"}`, http.StatusBadRequest)
		return
	}

	idStr := pathParts[len(pathParts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, `{"error": "Invalid task ID"}`, http.StatusBadRequest)
		return
	}

	// Маршрутизация по методу
	switch r.Method {
	case http.MethodGet:
		task, exists := h.Store.Get(id)
		if !exists {
			http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)

	case http.MethodPut, http.MethodPatch:
		// Обновить задачу (заглушка ШАГ 4)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":      id,
			"message": "Task updated (step 4)",
		})

	case http.MethodDelete:
		err := h.Store.Delete(id)
		if err != nil {
			http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

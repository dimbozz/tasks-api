// internal/handlers/tasks.go (заглушки обработчиков)
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"tasks-api/internal/models"
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
		// реализация POST /tasks
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, `{"error": "Content-Type must be application/json"}`, http.StatusUnsupportedMediaType)
			return
		}

		var task models.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
			return
		}

		// Валидация обязательных полей
		if task.Title == "" {
			http.Error(w, `{"error": "Title is required"}`, http.StatusBadRequest)
			return
		}

		// Создание задачи через Storage
		createdTask, err := h.Store.Create(task)
		if err != nil {
			http.Error(w, `{"error": "Failed to create task"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", fmt.Sprintf("/tasks/%d", createdTask.ID))
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdTask)

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
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, `{"error": "Content-Type must be application/json"}`, http.StatusUnsupportedMediaType)
			return
		}

		var update models.Task
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
			return
		}

		// Обновление через Storage.Update()
		updatedTask, err := h.Store.Update(id, update)
		if err != nil {
			http.Error(w, `{"error": "Task not found"}`, http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedTask)

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

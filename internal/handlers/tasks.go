// internal/handlers/tasks.go (заглушки обработчиков)
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"tasks-api/internal/models"
	"tasks-api/internal/storage"
)

type Handler struct{ Store storage.Storage }

func New(s storage.Storage) *Handler { return &Handler{Store: s} }

// Создаем единый формат ошибки JSON Error
type ErrorResponse struct {
	Error string `json:"error"`
}

func sendError(w http.ResponseWriter, status int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: errMsg})
}

func sendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// /tasks (GET, POST)
func (h *Handler) TasksCollection(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path) // Логирование запроса
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// Список задач
		tasks := h.Store.List()
		sendJSON(w, http.StatusOK, tasks)

	case http.MethodPost:
		// Проверка Content-Type
		if r.Header.Get("Content-Type") != "application/json" {
			sendError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
			return
		}

		var task models.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			sendError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		// Валидация обязательных полей
		if task.Title == "" {
			sendError(w, http.StatusBadRequest, "Title is required")
			return
		}

		// Создание задачи через Storage
		createdTask, err := h.Store.Create(task)
		if err != nil {
			sendError(w, http.StatusInternalServerError, "Failed to create task")
			return
		}

		w.Header().Set("Location", fmt.Sprintf("/tasks/%d", createdTask.ID))
		sendJSON(w, http.StatusCreated, createdTask)

	default:
		sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// /tasks/{id} (GET, PUT, DELETE)
func (h *Handler) TaskItem(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path) // Логирование запроса
	w.Header().Set("Content-Type", "application/json")

	// Извлечение ID из пути
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		sendError(w, http.StatusBadRequest, "Task ID required")
		return
	}

	idStr := pathParts[len(pathParts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		sendError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	// Маршрутизация по методу
	switch r.Method {
	case http.MethodGet:
		task, exists := h.Store.Get(id)
		if !exists {
			sendError(w, http.StatusNotFound, "Task not found")
			return
		}
		sendJSON(w, http.StatusOK, task)

	case http.MethodPut, http.MethodPatch:
		if r.Header.Get("Content-Type") != "application/json" {
			sendError(w, http.StatusUnsupportedMediaType, "Content-Type must be application/json")
			return
		}

		var update models.Task
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			sendError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		// Обновление через Storage.Update()
		updatedTask, err := h.Store.Update(id, update)
		if err != nil {
			sendError(w, http.StatusNotFound, "Task not found")
			return
		}
		sendJSON(w, http.StatusOK, updatedTask)

	case http.MethodDelete:
		err := h.Store.Delete(id)
		if err != nil {
			sendError(w, http.StatusNotFound, "Task not found")
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

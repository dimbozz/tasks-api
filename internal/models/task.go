// internal/models/task.go (модель)
package models

type Task struct {
	ID        int    `json:"id"`                   // генерируется на сервере
	Title     string `json:"title"`                // обязательное поле
	Done      bool   `json:"done"`                 // статус выполнения
	CreatedAt string `json:"created_at,omitempty"` // ISO8601, по желанию
}

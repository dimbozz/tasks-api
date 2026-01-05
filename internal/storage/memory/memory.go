// internal/storage/memory.go
// реализация in‑memory хранилища
package memory

import (
	"errors"
	"sync"
	"tasks-api/internal/models"
	"tasks-api/internal/storage" // интерфейс Storage
	"time"
)

// Создаем MemoryStorage — потокобезопасное in-memory хранилище задач
type MemoryStorage struct {
	tasks  []models.Task // список всех задач
	mu     sync.RWMutex  // RWMutex: RLock для чтения, Lock для записи
	nextID int           // следующий ID для новых задач
}

// New создает новое хранилище и возвращает интерфейс Storage (Dependency Injection)
func New() storage.Storage {
	return &MemoryStorage{nextID: 1}
}

// Реализация метода List для соответствия интерфейсу
// возвращает КОПИЮ всех задач (потокобезопасное чтение)
func (s *MemoryStorage) List() []models.Task {
	s.mu.RLock() // блокируем запись, разрешаем множественное чтение (read-heavy)
	defer s.mu.RUnlock() // снимаем блокировку при выходе из функции
	tasks := make([]models.Task, len(s.tasks)) // создаем копию
	copy(tasks, s.tasks)                       // копируем данные
	return tasks                               // возвращаем копию (без race condition)
}

// Реализация метода Create для соответствия интерфейсу
// создает новую задачу с автоинкрементным ID и текущей датой
func (s *MemoryStorage) Create(task models.Task) (models.Task, error) {
	s.mu.Lock() // Эксклюзивная блокировка для безопасной записи
	defer s.mu.Unlock() // снимаем блокировку при выходе из функции

	task.ID = s.nextID              // устанавливаем уникальный ID
	task.CreatedAt = time.Now()     // текущая дата создания
	s.tasks = append(s.tasks, task) // добавляем в срез
	s.nextID++                      // увеличиваем счетчик
	return task, nil
}

// Реализация метода Get для соответствия интерфейсу
// ищет задачу по ID (thread-safe чтение)
func (s *MemoryStorage) Get(id int) (models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, t := range s.tasks { // линейный поиск
		if t.ID == id {
			return t, true // найдена: возвращаем задачу + true
		}
	}
	return models.Task{}, false // не найдена: zero-value + false
}

// Реализация метода Update для соответствия интерфейсу
// обновляет существующую задачу, сохраняя дату создания
func (s *MemoryStorage) Update(id int, updTask models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, t := range s.tasks { // ищем по ID
		if t.ID == id {
			updTask.ID = id                 // сохраняем оригинальный ID
			updTask.CreatedAt = t.CreatedAt // НЕ меняем дату создания
			s.tasks[i] = updTask            // заменяем в срезе
			return updTask, nil
		}
	}
	return models.Task{}, errors.New("task not found") // не найдена
}

// Реализация метода Delete для соответствия интерфейсу
// удаляет задачу по ID
func (s *MemoryStorage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, t := range s.tasks { // ищем по ID
		if t.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...) // удаляем элемент
			return nil
		}
	}
	return errors.New("task not found") // не найдена
}

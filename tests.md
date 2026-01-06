Эндпоинты:
- GET /tasks — получить список задач.
```
curl http://localhost:8081/tasks
[{"id":1,"title":"Купить молоко","done":true,"created_at":"2026-01-06T16:07:33.576124+03:00"},{"id":2,"title":"Купить сахар","done":false,"created_at":"2026-01-06T16:16:39.91403+03:00"}]
```
- POST /tasks — создать задачу.
```
сurl -X POST http://localhost:8081/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "Купить сахар", "description": "20 кг", "done": false}'
{"id":2,"title":"Купить сахар","done":false,"created_at":"2026-01-06T16:16:39.91403+03:00"}

curl -X POST http://localhost:8081/tasks \
  -H "Content-Type: application/" \    
  -d '{"title": "Купить сахар", "description": "20 кг", "done": false}'
{"error":"Content-Type must be application/json"}
```
- GET /tasks/{id} — получить задачу по идентификатору.
```
curl http://localhost:8081/tasks/1
```
- PUT /tasks/{id} — обновить задачу целиком.
```
```
- DELETE /tasks/{id} — удалить задачу.
```
```
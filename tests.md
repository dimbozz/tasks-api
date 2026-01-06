Тесты по эндпоинтам:
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
curl http://localhost:8081/tasks/2
{"id":2,"title":"Купить сахар","done":false,"created_at":"2026-01-06T16:16:39.91403+03:00"}

curl -v http://localhost:8081/tasks/3
* Host localhost:8081 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:8081...
* Connected to localhost (::1) port 8081
> GET /tasks/3 HTTP/1.1
> Host: localhost:8081
> User-Agent: curl/8.7.1
> Accept: */*
> 
* Request completely sent off
< HTTP/1.1 404 Not Found
< Content-Type: application/json
< Date: Tue, 06 Jan 2026 13:21:31 GMT
< Content-Length: 27
< 
{"error":"Task not found"}
* Connection #0 to host localhost left intact
```
- PUT /tasks/{id} — обновить задачу целиком.
```
curl -X PATCH http://localhost:8081/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Купить молоко", "done": true}'
{"id":1,"title":"Купить молоко","done":true,"created_at":"2026-01-06T16:07:33.576124+03:00"}

curl -X PATCH http://localhost:8081/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": Купить молоко", "done": true}' 
{"error":"Invalid JSON"}
```
- DELETE /tasks/{id} — удалить задачу.

```
curl -X DELETE http://localhost:8081/tasks/2
[{"id":1,"title":"Купить молоко","done":true,"created_at":"2026-01-06T16:07:33.576124+03:00"}]

curl -X DELETE http://localhost:8081/tasks/3
{"error":"Task not found"}
```
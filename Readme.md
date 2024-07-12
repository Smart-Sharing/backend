# Sharing Warehouse Machines

**How to up postgres database**:
```
docker compose up -d
```

Run app:
```
go build cmd/main.go
./main --config=config/local.yml
```


## Frontend API Examples

`curl` - утилита для отправки запросов. Если нету `curl`, то отправить запрос можно любым другим способом.

```
# пример получения всех пользователей/машин/сессий
curl "localhost:8080/get_all_users" -X GET
curl "localhost:8080/get_all_machines" -X GET
curl "localhost:8080/get_all_sessions" -X GET

# примеры получения объектов по их id
curl "localhost:8080/get_user?user_id=1" -X GET
curl "localhost:8080/get_machine?machine_id=1FGH345" -X GET
curl "localhost:8080/get_session?session_id=2" -X GET
```

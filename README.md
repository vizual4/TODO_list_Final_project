# Файлы для итогового задания
Проект нужен для того, чтобы разгрузить память, не забывать важные задачи и эффективно планировать своё время. Он может добавлять задачи, переносить их на заданный интервал времени, искать их по ключевым словам, удалять.

Все задания со звездочкой были выполнены.

В директории `pkg/api` находятся хэндлеры для поисковых запросов `api.go, tasks.go`, а также хэндлеры для аутентификации `auth.go` и функция для вычисления следующей даты `nextdate.go`

В директории `pkg/db` находится файл для создания базы данных `db.go` и файл с функциями для взаимодействия с базой данных `task.go`

В директории `pgk/server` находится файл для запуска веб-сервера `server.go`

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

параметры для тестов следует указывать следующие:
package tests

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJoYXNoIjoiNTk5NDQ3MWFiYjAxMTEyYWZjYzE4MTU5ZjZjYzc0YjRmNTExYjk5ODA2ZGE1OWIzY2FmNWE5YzE3M2NhY2ZjNSJ9.c61UEKKEv6Ome3eD6owN--R7GQxXe3QmkHy9G-sJHLM`

ДЛЯ ЗАПУСКА ЛОКАЛЬНО:
Команда для запуска: TODO_PORT=7540 TODO_DBFILE=./scheduler.db TODO_PASSWORD=12345 go run main.go
Адрес для доступа в браузере: http://localhost:7540

ДЛЯ ЗАПУСКА ИЗ КОНТЕЙНЕРА:
Сборка Docker-образа производится командой `docker build -t todo-app:v1 .`

Для запуска контейнера необходимо ввести следующую команду: 
docker run -d \
  --name todo-container \
  -p 7540:7540 \
  -e TODO_PASSWORD="123" \
  -e TODO_PORT="7540" \
  -e TODO_DBFILE="/app/db/scheduler.db" \
  -v "$(pwd)/db:/app/db" \
  todo-app:v1
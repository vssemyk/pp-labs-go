# Лабораторная №8 - REST API на Go + PostgreSQL

## Выполненные задания

1. Построен базовый REST API на стандартном пакете `net/http` с поддержкой маршрутов `GET /users`, `GET /users/{id}`, `POST /users`, `PUT /users/{id}` и `DELETE /users/{id}`.
2. Подключена реляционная база данных PostgreSQL (`database/sql` + драйвер `github.com/lib/pq`) для хранения и обработки записей пользователей.
3. Реализована централизованная обработка ошибок в едином JSON-формате (`sendError`) и валидация входных данных (`validateUser`) при создании и обновлении записей.
4. Добавлена поддержка пагинации (`page`, `limit`) и фильтрации (`name` без учёта регистра через `ILIKE` и `age` по точному совпадению) в маршруте `GET /users`.
5. Написаны модульные тесты (`handlers_test.go`) с использованием пакета `net/http/httptest` для проверки валидации и обработки некорректных запросов.
6. Подготовлена полная документация API с описанием эндпоинтов, моделей данных, кодов ошибок и примеров `curl`-запросов.

## Требования

- Go 1.20+
- PostgreSQL 11+
- Драйвер `github.com/lib/pq`

## Установка и запуск

1. Создайте базу данных и таблицу в PostgreSQL:

```sql
CREATE DATABASE usersdb;
\c usersdb

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    age INT NOT NULL
);
```

2. Инициализируйте модуль, установите драйвер PostgreSQL и проверьте строку подключения в `db.go`:

```bash
go mod init users-api
go get github.com/lib/pq
```

3. Запустите сервер из директории проекта:

```bash
go run .
```

При корректном запуске в консоль выводится:
```
Подключение к PostgreSQL успешно
Сервер запущен на :8080
```

Для запуска unit-тестов выполните:

```bash
go test -v ./...
```

## Структура проекта

| Файл | Назначение |
|---|---|
| `main.go` | Регистрация маршрутов и запуск HTTP-сервера |
| `db.go` | Инициализация и проверка подключения к PostgreSQL |
| `model.go` | Описание структуры `User` с JSON-тегами |
| `handlers.go` | Обработчики CRUD-маршрутов, пагинации и фильтрации |
| `validate.go` | Функция валидации полей пользователя (`name`, `email`, `age`) |
| `errors.go` | Централизованная отправка JSON-ошибок (`APIError`, `sendError`) |
| `handlers_test.go` | Unit-тесты для валидации и обработчиков |
| `README.md` | Документация REST API |

## Задание 1 и 2 - Построение базового REST API и подключение PostgreSQL
Сервер запускается на порту `:8080` и распределяет запросы между маршрутами `/users` и `/users/{id}` с проверкой допустимых HTTP-методов. Взаимодействие с базой данных `usersdb` реализовано через параметризованные SQL-запросы (`$1`, `$2`, ...), что защищает от SQL-инъекций.

### Модель User

```json
{
  "id": 1,
  "name": "Влада",
  "email": "solov@mail.ru",
  "age": 26
}
```

### Пример - Создание пользователя (`POST /users`)
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Влада\",\"email\":\"solov@mail.ru\",\"age\":26}"
```

### Ответ
`201 Created`
```json
{"id": 1, "name": "Влада", "email": "solov@mail.ru", "age": 26}
```

### Пример - Получение пользователя по ID (`GET /users/{id}`)
`GET /users/1`

### Ответ
`200 OK`
```json
{"id": 1, "name": "Влада", "email": "solov@mail.ru", "age": 26}
```

### Пример - Обновление пользователя (`PUT /users/{id}`)
```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Влада Соловова\",\"email\":\"vsolov@mail.ru\",\"age\":27}"
```

### Ответ
`200 OK`
```json
{"id": 1, "name": "Влада Соловова", "email": "vsolov@mail.ru", "age": 27}
```

### Пример - Удаление пользователя (`DELETE /users/{id}`)
```bash
curl -X DELETE http://localhost:8080/users/1
```

### Ответ
`204 No Content`

## Задание 3 - Обработка ошибок и валидация данных
Все ошибки возвращаются в едином формате `{"code": "<код>", "message": "<описание>"}`. Перед сохранением в БД функция `validateUser` проверяет, что имя не пустое, `email` содержит символ `@`, а возраст находится в диапазоне от `0` до `150`.

### Таблица кодов ошибок

| HTTP | code | Когда возникает |
|---|---|---|
| 400 | `bad_json` | Тело запроса не является корректным JSON |
| 400 | `bad_id` | ID в адресе не является целым числом |
| 400 | `validation` | Данные не прошли валидацию |
| 404 | `not_found` | Пользователь с таким ID не найден |
| 405 | `method` | HTTP-метод не поддерживается для маршрута |
| 500 | `db_error` | Ошибка при выполнении запроса к БД |

### Тест - Некорректный JSON
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d "не json"
```

### Ответ
`400 Bad Request`
```json
{"code": "bad_json", "message": "Некорректный JSON"}
```

### Тест - Ошибка валидации (пустое имя, плохой email или возраст вне диапазона)
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"\",\"email\":\"bad\",\"age\":-1}"
```

### Ответ
`400 Bad Request`
```json
{"code": "validation", "message": "Имя не может быть пустым"}
```

### Тест - Некорректный ID в URL
`GET /users/abc`

### Ответ
`400 Bad Request`
```json
{"code": "bad_id", "message": "Некорректный ID"}
```

### Тест - Пользователь не найден
`GET /users/999`

### Ответ
`404 Not Found`
```json
{"code": "not_found", "message": "Пользователь не найден"}
```

### Тест - Неподдерживаемый HTTP-метод
`PATCH /users`

### Ответ
`405 Method Not Allowed`
```json
{"code": "method", "message": "Метод не разрешён"}
```

## Задание 4 - Пагинация и фильтрация
Маршрут `GET /users` поддерживает необязательные query-параметры: `page` (по умолчанию `1`), `limit` (по умолчанию `10`), `name` (поиск подстроки без учёта регистра) и `age` (точное совпадение).

### Пример
```
GET /users?name=Влада&age=26&page=1&limit=5
```

### Ответ
`200 OK`
```json
[
  {"id": 1, "name": "Влада", "email": "solov@mail.ru", "age": 26}
]
```

## Задание 5 - Тестирование API
В файле `handlers_test.go` реализованы табличные unit-тесты `TestValidateUser` (проверка валидного пользователя, пустого имени, некорректного email, отрицательного возраста) и интеграционный тест обработчика `TestCreateUserBadJSON` с использованием `httptest`.

### Пример
`go test -v ./...`

### Ответ
```
=== RUN   TestValidateUser
=== RUN   TestValidateUser/валидный
=== RUN   TestValidateUser/пустое_имя
=== RUN   TestValidateUser/плохой_email
=== RUN   TestValidateUser/отрицательный_возраст
--- PASS: TestValidateUser (0.00s)
=== RUN   TestCreateUserBadJSON
--- PASS: TestCreateUserBadJSON (0.00s)
PASS
ok      users-api   0.012s
```

## Задание 6 - Документация API
Полная спецификация эндпоинтов, форматов запросов и ответов сведена в таблицу маршрутов:

| Метод | Путь | Описание |
|---|---|---|
| GET | `/users` | Список пользователей (параметры `page`, `limit`, `name`, `age`) |
| GET | `/users/{id}` | Получение пользователя по ID |
| POST | `/users` | Создание нового пользователя |
| PUT | `/users/{id}` | Полное обновление данных пользователя |
| DELETE | `/users/{id}` | Удаление пользователя по ID |

## Исходный код

### `main.go`
```go
package main

import (
	"log"
	"net/http"
)

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getUsers(w, r)
		case http.MethodPost:
			createUser(w, r)
		default:
			sendError(w, http.StatusMethodNotAllowed, "method", "Метод не разрешён")
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getUser(w, r)
		case http.MethodPut:
			updateUser(w, r)
		case http.MethodDelete:
			deleteUser(w, r)
		default:
			sendError(w, http.StatusMethodNotAllowed, "method", "Метод не разрешён")
		}
	})

	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### `db.go`
```go
package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	connStr := "user=postgres password=2562 dbname=usersdb sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка открытия БД:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("БД недоступна:", err)
	}

	log.Println("Подключение к PostgreSQL успешно")
}
```

### `model.go`
```go
package main

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}
```

### `handlers.go`
```go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func getUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	ageStr := query.Get("age")
	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	page := 1
	limit := 10
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	offset := (page - 1) * limit

	sqlQuery := "SELECT id, name, email, age FROM users WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if name != "" {
		sqlQuery += fmt.Sprintf(" AND name ILIKE $%d", argNum)
		args = append(args, "%"+name+"%")
		argNum++
	}
	if ageStr != "" {
		if age, err := strconv.Atoi(ageStr); err == nil {
			sqlQuery += fmt.Sprintf(" AND age = $%d", argNum)
			args = append(args, age)
			argNum++
		}
	}

	sqlQuery += fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, limit, offset)

	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "db_error", err.Error())
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age)
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "bad_id", "Некорректный ID")
		return
	}

	var u User
	err = db.QueryRow("SELECT id, name, email, age FROM users WHERE id = $1", id).
		Scan(&u.ID, &u.Name, &u.Email, &u.Age)

	if err == sql.ErrNoRows {
		sendError(w, http.StatusNotFound, "not_found", "Пользователь не найден")
		return
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, "db_error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		sendError(w, http.StatusBadRequest, "bad_json", "Некорректный JSON")
		return
	}

	if msg := validateUser(u); msg != "" {
		sendError(w, http.StatusBadRequest, "validation", msg)
		return
	}

	err := db.QueryRow(
		"INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id",
		u.Name, u.Email, u.Age,
	).Scan(&u.ID)

	if err != nil {
		sendError(w, http.StatusInternalServerError, "db_error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "bad_id", "Некорректный ID")
		return
	}

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		sendError(w, http.StatusBadRequest, "bad_json", "Некорректный JSON")
		return
	}

	if msg := validateUser(u); msg != "" {
		sendError(w, http.StatusBadRequest, "validation", msg)
		return
	}

	result, err := db.Exec(
		"UPDATE users SET name=$1, email=$2, age=$3 WHERE id=$4",
		u.Name, u.Email, u.Age, id,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "db_error", err.Error())
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		sendError(w, http.StatusNotFound, "not_found", "Пользователь не найден")
		return
	}

	u.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "bad_id", "Некорректный ID")
		return
	}

	result, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "db_error", err.Error())
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		sendError(w, http.StatusNotFound, "not_found", "Пользователь не найден")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
```

### `validate.go`
```go
package main

import (
	"strings"
)

func validateUser(u User) string {
	if strings.TrimSpace(u.Name) == "" {
		return "Имя не может быть пустым"
	}
	if !strings.Contains(u.Email, "@") {
		return "Некорректный email"
	}
	if u.Age < 0 || u.Age > 150 {
		return "Возраст должен быть от 0 до 150"
	}
	return ""
}
```

### `errors.go`
```go
package main

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func sendError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIError{
		Code:    code,
		Message: message,
	})
}
```

### `handlers_test.go`
```go
package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateUser(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{"валидный", User{Name: "Иван", Email: "ivan@mail.ru", Age: 25}, false},
		{"пустое имя", User{Name: "", Email: "a@b.ru", Age: 20}, true},
		{"плохой email", User{Name: "A", Email: "bad", Age: 20}, true},
		{"отрицательный возраст", User{Name: "A", Email: "a@b.ru", Age: -1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := validateUser(tt.user)
			if (msg != "") != tt.wantErr {
				t.Errorf("validateUser() = %q, wantErr %v", msg, tt.wantErr)
			}
		})
	}
}

func TestCreateUserBadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader("не json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	createUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("код = %d, ожидался %d", w.Code, http.StatusBadRequest)
	}
}
```
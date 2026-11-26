# Лабораторная №9 - Клиент-серверное приложение с авторизацией на Go

## Выполненные задания

1. Разработано консольное клиентское приложение на Go, отправляющее HTTP-запросы к REST API серверу управления пользователями.
2. Реализован механизм обработки HTTP-ответов и ошибок в едином формате JSON, а также функции форматированного табличного вывода данных.
3. Создан интерактивный консольный интерфейс с меню для выполнения всех CRUD-операций (просмотр списка с фильтрацией, получение по ID, создание, обновление и удаление пользователей).
4. Добавлена авторизация по токену: сервер генерирует сессионный токен при входе через `POST /login`, а клиент сохраняет его и передаёт в заголовке `Authorization: Bearer <token>`.
5. Обеспечена одновременная работа нескольких клиентов с сервером за счёт потокобезопасного хранилища токенов с использованием `sync.Mutex`.
6. Проведено тестирование взаимодействия клиента и сервера, а также написаны unit-тесты для проверки валидации и обработки некорректных запросов.

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

2. Добавьте хотя бы одного пользователя в таблицу `users` для первичного входа через `/login`, например:

```sql
INSERT INTO users (name, email, age) VALUES ('Влада', 'solov@mail.ru', 26);
```

3. Перейдите в директорию проекта `lb9/users-api`, установите зависимости и запустите сервер:

```bash
go mod tidy
go run .
```

Сервер запустится по адресу `http://localhost:8080`.

4. В отдельном терминале (или нескольких терминалах одновременно) запустите консольный клиент:

```bash
go run ./client
```

После запуска клиента откроется интерактивное меню с выбором действий от `0` до `6`. Без предварительного входа (`1 - Войти`) запросы к `/users` вернут ошибку `401 Unauthorized`.

## Структура проекта

| Файл | Назначение |
|---|---|
| `client/main.go` | Консольное клиентское приложение с меню, хранением токена и выполнением CRUD-запросов |
| `main.go` | Точка входа сервера, регистрация маршрутов `/login`, `/users`, `/users/` и проверка авторизации |
| `auth.go` | Генерация токенов сессии, обработчик `POST /login` и middleware-проверка `Authorization` |
| `db.go` | Инициализация подключения к базе данных PostgreSQL |
| `model.go` | Описание структуры данных `User` |
| `handlers.go` | Обработчики CRUD-операций, пагинации и фильтрации |
| `validate.go` | Валидация полей пользователя при создании и обновлении |
| `errors.go` | Централизованная отправка JSON-ошибок в едином формате |
| `handlers_test.go` | Unit-тесты валидации и обработки некорректного JSON |
| `go.mod` | Описание модуля Go и зависимостей |
| `go.sum` | Контрольные суммы зависимостей |

## Задание 1 - Создание клиентской программы
Разработан HTTP-клиент в файле `client/main.go`. Функция `doRequest` формирует HTTP-запросы (`GET`, `POST`, `PUT`, `DELETE`) к базовому адресу `http://localhost:8080`, сериализует тело запроса в JSON, устанавливает заголовок `Content-Type: application/json` и таймаут ожидания `5 * time.Second`.

### Пример
```
=== Клиент Users API ===
Сервер: http://localhost:8080
```

### Ответ
Клиент устанавливает соединение с сервером при выборе любого пункта меню и выводит статус выполнения операции.

## Задание 2 - Обработка ответов
Реализована функция `printError` для разбора JSON-ответов сервера со структурой `APIError` (`code`, `message`), а также функции `printUser` и `printUsers` для аккуратного форматированного вывода записей в консоль.

### Пример
Запрос списка пользователей `GET /users` через пункт меню `2`.

### Ответ
```
Список пользователей:
  id=1   name="Влада"              email="solov@mail.ru"           age=26
  id=2   name="Иван"               email="ivan@mail.ru"            age=25
```

### Тест - Сервер недоступен
Запуск запроса при выключенном сервере.

### Ответ
`❌ Не удалось соединиться с сервером: Post "http://localhost:8080/login": dial tcp [::1]:8080: connectex: No connection could be made because the target machine actively refused it.`

## Задание 3 - Интерфейс пользователя
В функции `main()` клиента реализован бесконечный цикл с текстовым меню для выполнения всех CRUD-операций: просмотр списка с опциональной фильтрацией по имени, поиск пользователя по ID, создание, обновление и удаление записи.

### Пример
```
Меню:
  1 - Войти (получить токен)
  2 - Показать список пользователей
  3 - Показать пользователя по id
  4 - Создать пользователя
  5 - Обновить пользователя
  6 - Удалить пользователя
  0 - Выйти из программы
Выбор: 4
Имя: Эвелина
Email: ex@mail.ru
Возраст: 30
```

### Ответ
```
✅ Создан:
  id=3   name="Эвелина"            email="ex@mail.ru"            age=30
```

### Тест - Некорректный ввод возраста или ID
```
Выбор: 3
ID пользователя: abc
```

### Ответ
`Некорректный ID`

### Тест - Ошибка валидации при создании пользователя
```
Выбор: 4
Имя: Эвелина
Email: invalid-email
Возраст: 30
```

### Ответ
`❌ Ошибка 400 [validation]: Некорректный email`

## Задание 4 - Авторизация пользователя
На сервере в файле `auth.go` добавлен маршрут `POST /login`, который проверяет наличие пользователя по `email` в базе данных и возвращает сгенерированный 32-символьный hex-токен. Клиент сохраняет токен в глобальной переменной `sessionToken` и автоматически подставляет заголовок `Authorization: Bearer <token>` во все последующие запросы.

### Пример
```
Выбор: 1
Email: solov@mail.ru
```

### Ответ
`✅ Вы вошли. Токен сохранён.`

### Тест - Запрос без авторизации
Попытка вызвать любую CRUD-операцию (`2-6`) до выполнения входа.

### Ответ
`⚠️  Вы не авторизованы. Сначала выполните вход.`

### Тест - Вход с несуществующим email
```
Выбор: 1
Email: unknown@mail.ru
```

### Ответ
`❌ Ошибка 401 [auth]: Пользователь с таким email не найден`

### Тест - Запрос к серверу без заголовка Authorization
Выполнение прямого HTTP-запроса `GET /users` без токена.

### Ответ
`{"code":"auth","message":"Требуется авторизация"}`

## Задание 5 - Поддержка нескольких клиентов
На сервере карта активных сессий `tokens = map[string]int{}` защищена мьютексом `tokensMu sync.Mutex` при чтении и записи. Это исключает состояние гонки (data race) при одновременном входе и выполнении запросов от нескольких параллельно запущенных экземпляров клиента.

### Пример
Запуск двух экземпляров `go run ./client` в разных окнах терминала с авторизацией под разными или одинаковыми учетными записями.

### Ответ
Оба клиента получают независимые токены сессии и корректно выполняют чтение и изменение данных в общей базе PostgreSQL без конфликтов.

## Задание 6 - Тестирование и отладка
Проверены все маршруты (`POST /login`, `GET /users`, `GET /users/{id}`, `POST /users`, `PUT /users/{id}`, `DELETE /users/{id}`), фильтрация по имени, обработка кодов ошибок (`400`, `401`, `404`, `405`, `500`), а также запущены модульные тесты `go test -v ./...`.

### Пример
```bash
go test -v ./...
```

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
ok      users-api       0.215s
```

### Тест - Запрос несуществующего пользователя
```
Выбор: 3
ID пользователя: 999
```

### Ответ
`❌ Ошибка 404 [not_found]: Пользователь не найден`

## Исходный код

### `client/main.go`
```go
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const baseURL = "http://localhost:8080"

var sessionToken string

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

var reader = bufio.NewReader(os.Stdin)

func readLine(prompt string) string {
	fmt.Print(prompt)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func readInt(prompt string) (int, bool) {
	s := readLine(prompt)
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return n, true
}

func doRequest(method, path string, body interface{}) (int, []byte, error) {
	var reqBody io.Reader

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return 0, nil, err
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, baseURL+path, reqBody)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	if sessionToken != "" {
		req.Header.Set("Authorization", "Bearer "+sessionToken)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, respBody, nil
}

func printError(status int, body []byte) {
	var e APIError
	if err := json.Unmarshal(body, &e); err == nil && e.Code != "" {
		fmt.Printf("❌ Ошибка %d [%s]: %s\n", status, e.Code, e.Message)
	} else {
		fmt.Printf("❌ Ошибка %d: %s\n", status, string(body))
	}
}

func printUser(u User) {
	fmt.Printf("  id=%-3d name=%-20q email=%-25q age=%d\n", u.ID, u.Name, u.Email, u.Age)
}

func printUsers(users []User) {
	if len(users) == 0 {
		fmt.Println("  (список пуст)")
		return
	}
	for _, u := range users {
		printUser(u)
	}
}

// 1. Логин
func doLogin() {
	email := readLine("Email: ")
	if email == "" {
		fmt.Println("Email не может быть пустым")
		return
	}

	status, body, err := doRequest("POST", "/login", map[string]string{"email": email})
	if err != nil {
		fmt.Println("❌ Не удалось соединиться с сервером:", err)
		return
	}
	if status != http.StatusOK {
		printError(status, body)
		return
	}

	var resp struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Println("❌ Не удалось разобрать ответ:", err)
		return
	}

	sessionToken = resp.Token
	fmt.Println("✅ Вы вошли. Токен сохранён.")
}

func isLoggedIn() bool {
	return sessionToken != ""
}

// 2. Показать список пользователей
func doList() {

	if !isLoggedIn() {
		fmt.Println("⚠️  Вы не авторизованы. Сначала выполните вход.")
		return
	}

	path := "/users"

	name := readLine("Фильтр по имени (Enter — пропустить): ")
	if name != "" {
		path += "?name=" + name
	}

	status, body, err := doRequest("GET", path, nil)
	if err != nil {
		fmt.Println("❌ Не удалось соединиться с сервером:", err)
		return
	}
	if status != http.StatusOK {
		printError(status, body)
		return
	}

	var users []User
	if err := json.Unmarshal(body, &users); err != nil {
		fmt.Println("❌ Не удалось разобрать ответ:", err)
		return
	}

	fmt.Println("Список пользователей:")
	printUsers(users)
}

// 3. Показать одного пользователя
func doGet() {

	if !isLoggedIn() {
		fmt.Println("⚠️  Вы не авторизованы. Сначала выполните пункт вход.")
		return
	}

	id, ok := readInt("ID пользователя: ")
	if !ok {
		fmt.Println("Некорректный ID")
		return
	}

	status, body, err := doRequest("GET", "/users/"+strconv.Itoa(id), nil)
	if err != nil {
		fmt.Println("❌ Не удалось соединиться с сервером:", err)
		return
	}
	if status != http.StatusOK {
		printError(status, body)
		return
	}

	var u User
	if err := json.Unmarshal(body, &u); err != nil {
		fmt.Println("❌ Не удалось разобрать ответ:", err)
		return
	}
	fmt.Println("Пользователь:")
	printUser(u)
}

// 4. Создать пользователя
func doCreate() {

	if !isLoggedIn() {
		fmt.Println("⚠️  Вы не авторизованы. Сначала выполните пункт вход.")
		return
	}

	name := readLine("Имя: ")
	email := readLine("Email: ")
	age, ok := readInt("Возраст: ")
	if !ok {
		fmt.Println("Некорректный возраст")
		return
	}

	u := User{Name: name, Email: email, Age: age}

	status, body, err := doRequest("POST", "/users", u)
	if err != nil {
		fmt.Println("❌ Не удалось соединиться с сервером:", err)
		return
	}
	if status != http.StatusCreated {
		printError(status, body)
		return
	}

	var created User
	json.Unmarshal(body, &created)
	fmt.Println("✅ Создан:")
	printUser(created)
}

// 5. Обновить пользователя
func doUpdate() {

	if !isLoggedIn() {
		fmt.Println("⚠️  Вы не авторизованы. Сначала выполните пункт вход.")
		return
	}

	id, ok := readInt("ID пользователя для обновления: ")
	if !ok {
		fmt.Println("Некорректный ID")
		return
	}

	name := readLine("Новое имя: ")
	email := readLine("Новый email: ")
	age, ok := readInt("Новый возраст: ")
	if !ok {
		fmt.Println("Некорректный возраст")
		return
	}

	u := User{ID: id, Name: name, Email: email, Age: age}

	status, body, err := doRequest("PUT", "/users/"+strconv.Itoa(id), u)
	if err != nil {
		fmt.Println("❌ Не удалось соединиться с сервером:", err)
		return
	}
	if status != http.StatusOK {
		printError(status, body)
		return
	}

	var updated User
	json.Unmarshal(body, &updated)
	fmt.Println("✅ Обновлён:")
	printUser(updated)
}

// 6. Удалить пользователя
func doDelete() {

	if !isLoggedIn() {
		fmt.Println("⚠️  Вы не авторизованы. Сначала выполните пункт вход.")
		return
	}

	id, ok := readInt("ID пользователя для удаления: ")
	if !ok {
		fmt.Println("Некорректный ID")
		return
	}

	status, body, err := doRequest("DELETE", "/users/"+strconv.Itoa(id), nil)
	if err != nil {
		fmt.Println("ХХХ Не удалось соединиться с сервером:", err)
		return
	}
	if status != http.StatusNoContent {
		printError(status, body)
		return
	}

	fmt.Println("✅ Удалено.")
}

func main() {
	fmt.Println("=== Клиент Users API ===")
	fmt.Println("Сервер:", baseURL)

	for {
		fmt.Println()
		fmt.Println("Меню:")
		fmt.Println("  1 - Войти (получить токен)")
		fmt.Println("  2 - Показать список пользователей")
		fmt.Println("  3 - Показать пользователя по id")
		fmt.Println("  4 - Создать пользователя")
		fmt.Println("  5 - Обновить пользователя")
		fmt.Println("  6 - Удалить пользователя")
		fmt.Println("  0 - Выйти из программы")

		if sessionToken == "" {
			fmt.Println("(вы не авторизованы — операции 2-6 вернут 401)")
		}

		choice := readLine("Выбор: ")

		switch choice {
		case "1":
			doLogin()
		case "2":
			doList()
		case "3":
			doGet()
		case "4":
			doCreate()
		case "5":
			doUpdate()
		case "6":
			doDelete()
		case "0":
			fmt.Println("Выход!")
			return
		default:
			fmt.Println("Не понял команду, попробуй ещё раз.")
		}
	}
}
```

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

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			sendError(w, http.StatusMethodNotAllowed, "method", "Метод не разрешён")
			return
		}
		loginHandler(w, r)
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if !requireAuth(w, r) {
			return
		}
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
		if !requireAuth(w, r) {
			return
		}
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

### `auth.go`
```go
package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
)

var (
	tokensMu sync.Mutex
	tokens   = map[string]int{}
)

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "bad_json", "Некорректный JSON")
		return
	}
	if req.Email == "" {
		sendError(w, http.StatusBadRequest, "validation", "Email обязателен")
		return
	}

	var id int
	err := db.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&id)
	if err != nil {
		sendError(w, http.StatusUnauthorized, "auth", "Пользователь с таким email не найден")
		return
	}

	token := generateToken()

	tokensMu.Lock()
	tokens[token] = id
	tokensMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func checkAuth(r *http.Request) (int, bool) {
	h := r.Header.Get("Authorization")
	if len(h) < 8 || h[:7] != "Bearer " {
		return 0, false
	}
	token := h[7:]

	tokensMu.Lock()
	id, ok := tokens[token]
	tokensMu.Unlock()

	return id, ok
}

func requireAuth(w http.ResponseWriter, r *http.Request) bool {
	if _, ok := checkAuth(r); !ok {
		sendError(w, http.StatusUnauthorized, "auth", "Требуется авторизация")
		return false
	}
	return true
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
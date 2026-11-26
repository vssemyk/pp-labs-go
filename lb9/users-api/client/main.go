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

# Лабораторная №10 - Криптография и безопасность приложений на Go

## Выполненные задания
- Задание 1: реализована утилита вычисления хэшей (`MD5`, `SHA-256`, `SHA-512`) и проверки целостности данных.
- Задание 2: реализовано симметричное шифрование и расшифрование строк по алгоритму `AES-GCM`.
- Задание 3: реализована генерация пары ключей `RSA` (`private.pem`, `public.pem`), подписание сообщений и проверка цифровой подписи.
- Задание 4: разработан защищённый TLS TCP-сервер и клиент с взаимной аутентификацией по сертификатам (mTLS).
- Задание 5: добавлена JWT-аутентификация в REST API и ролевое ограничение доступа (`admin`, `user`).
- Задание 6: внедрены защита от CSRF-атак через cookie и заголовок `X-CSRF-Token`, а также управление серверными сессиями.

## Требования
- Go `1.21` или выше
- База данных `PostgreSQL` (база `usersdb` с таблицей `users`)
- Зависимости из `go.mod`: `github.com/golang-jwt/jwt/v5`, `github.com/lib/pq`

## Установка и запуск
Сохраните все файлы проекта согласно структуре каталогов и установите зависимости в `users-api`:
```bash
cd users-api
go mod tidy
```

Запуск криптографической утилиты (задания 1–3):
```bash
cd users-api/crypto-lab
go run main.go
```
При корректном запуске в консоли появляется меню выбора криптографических операций (`1 - Вычислить хэш`, `3 - AES шифрование`, `4 - RSA: сгенерировать ключи` и др.).

Запуск TLS-сервера и клиента (задание 4):
```bash
cd users-api/tlslab
go run cert/main.go
go run server/main.go
# В отдельном терминале:
go run client/main.go
```
При корректном запуске создаются сертификаты `ca.crt`, `server.crt`, `client.crt`, сервер выводит `TLS TCP сервер запущен на :8443`, а клиент устанавливает взаимно аутентифицированное соединение.

Запуск защищённого REST API и консольного клиента (задания 5–6):
```bash
cd users-api
go run .
# В отдельном терминале:
go run client/main.go
```
При корректном запуске сервер выводит `Подключение к PostgreSQL успешно` и `Сервер запущен на :8080`.

## Структура проекта
| Файл | Назначение |
|---|---|
| `users-api/main.go` | Точка входа REST API, маршрутизация `/login`, `/users`, `/session/*`, `/csrf-demo` |
| `users-api/auth.go` | Выпуск и проверка JWT-токенов, проверка ролей (`user` и `admin`) |
| `users-api/csrf.go` | Middleware защиты от CSRF-атак и потокобезопасное хранилище сессий в памяти |
| `users-api/handlers.go` | CRUD-обработчики пользователей с фильтрацией и пагинацией |
| `users-api/handlers_test.go` | Модульные тесты валидации и обработки некорректного JSON |
| `users-api/db.go` | Подключение к PostgreSQL (`usersdb`) |
| `users-api/errors.go` | Единый формат JSON-ответов об ошибках API |
| `users-api/go.mod` | Описание модуля Go и внешних зависимостей |
| `users-api/go.sum` | Контрольные суммы зависимостей модуля |
| `users-api/client/main.go` | Консольный HTTP-клиент для работы с защищённым Users API по JWT |
| `users-api/crypto-lab/main.go` | Консольная утилита для хэширования, AES-шифрования и RSA-подписей |
| `users-api/tlslab/cert/main.go` | Генератор корневого CA, серверного и клиентского X.509-сертификатов |
| `users-api/tlslab/server/main.go` | Многопоточный TLS TCP-сервер с обязательной проверкой клиентского сертификата |
| `users-api/tlslab/client/main.go` | TLS TCP-клиент с проверкой сертификата сервера и отправкой своего сертификата |

## Задание 1 - Хэширование данных
Вычисление хэша строки с помощью `MD5`, `SHA-256` или `SHA-512` и проверка целостности данных путём сравнения фактического и ожидаемого значений хэша.

### Пример
```text
Выбор: 1
Строка: hello
Алгоритм (md5/sha256/sha512): sha256

Выбор: 2
Строка: hello
Ожидаемый хэш: 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824
Алгоритм (md5/sha256/sha512): sha256
```

### Ответ
`Хэш (sha256): 2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824`
`Хэши совпадают — данные целостны`

## Задание 2 - Симметричное шифрование
Шифрование введённой строки алгоритмом `AES-GCM` со случайным `nonce` и её расшифровка при повторном вводе секретного ключа (16, 24 или 32 байта).

### Пример
```text
Выбор: 3
Строка: Секретное сообщение
Ключ (16, 24 или 32 символа): 1234567890123456
Тот же ключ для расшифровки: 1234567890123456
```

### Ответ
`Зашифровано: 8a4f1c9e0b3d7a2f41c8e90b5d12a7f3...`
`Расшифровано: Секретное сообщение`

## Задание 3 - Асимметричное шифрование и цифровая подпись
Генерация пары ключей `RSA-2048` (`private.pem` и `public.pem`), подписание хэша `SHA-256` сообщения закрытым ключом и проверка подписи открытым ключом.

### Пример
```text
Выбор: 4
Выбор: 5
Сообщение: Договор №10
Выбор: 6
Сообщение: Договор №10
Подпись (hex): 4b82e9f1a0...
```

### Ответ
`Ключи сохранены: private.pem и public.pem`
`Подпись (hex): 4b82e9f1a0...`
`Подпись верна`

## Задание 4 - Реализация защищённого канала передачи данных (TLS)
Выпуск сертификатов удостоверяющего центра (`ca.crt`), сервера (`server.crt`) и клиента (`client.crt`), настройка TLS TCP-сервера на `:8443` и клиента с взаимной аутентификацией по сертификатам.

### Пример
```text
Введите сообщение: Привет по TLS!
```

### Ответ
`Подключено по TLS, сертификат сервера проверен через ca.crt`
`Клиентский сертификат предъявлен серверу`
`Ответ от сервера: Сообщение получено по TLS!`

## Задание 5 - Защита REST API
Выдача JWT-токена в `POST /login`, проверка заголовка `Authorization: Bearer <token>` на защищённых маршрутах `/users` и ограничение удаления (`DELETE /users/{id}`) только для роли `admin` (`id == 4`).

### Пример
```text
Выбор: 1
Email: admin@mail.ru
Выбор: 6
ID пользователя для удаления: 2
```

### Ответ
`✅ Вы вошли. Токен сохранён.`
`✅ Удалено.`

## Задание 6 - Защита от CSRF и управление сессиями
Управление серверными сессиями через `HttpOnly` cookie `session_id` (`/session/login`, `/session/whoami`) и защита от CSRF-атак через сверку cookie `csrf_token` и заголовка `X-CSRF-Token` (`/csrf-demo`).

### Пример
```bash
curl -i -c cookies.txt -X POST -d "email=ivan@mail.ru" http://localhost:8080/session/login
curl -i -b cookies.txt http://localhost:8080/session/whoami
curl -i -c csrf.txt http://localhost:8080/csrf-demo
curl -i -b csrf.txt -H "X-CSRF-Token: <token>" -X POST http://localhost:8080/csrf-demo
```

### Ответ
`Сессия создана`
`Ваш user_id: 1`
`CSRF-проверка пройдена`

## Исходный код

### `users-api/client/main.go`
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
	line, _ := reader.ReadString('
')
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
		fmt.Printf("❌ Ошибка %d [%s]: %s
", status, e.Code, e.Message)
	} else {
		fmt.Printf("❌ Ошибка %d: %s
", status, string(body))
	}
}
func printUser(u User) {
	fmt.Printf("  id=%-3d name=%-20q email=%-25q age=%d
", u.ID, u.Name, u.Email, u.Age)
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

### `users-api/crypto-lab/main.go`
```go
package main
import (
	"bufio"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"strings"
)
var reader = bufio.NewReader(os.Stdin)
func readLine(prompt string) string {
	fmt.Print(prompt)
	line, _ := reader.ReadString('
')
	return strings.TrimSpace(line)
}
func computeHash(data []byte, algo string) string {
	switch algo {
	case "md5":
		h := md5.Sum(data)
		return hex.EncodeToString(h[:])
	case "sha256":
		h := sha256.Sum256(data)
		return hex.EncodeToString(h[:])
	case "sha512":
		h := sha512.Sum512(data)
		return hex.EncodeToString(h[:])
	}
	return ""
}
func doHash() {
	text := readLine("Строка: ")
	algo := readLine("Алгоритм (md5/sha256/sha512): ")
	result := computeHash([]byte(text), algo)
	if result == "" {
		fmt.Println("Неизвестный алгоритм")
		return
	}
	fmt.Printf("Хэш (%s): %s
", algo, result)
}
func doVerify() {
	text := readLine("Строка: ")
	expected := readLine("Ожидаемый хэш: ")
	algo := readLine("Алгоритм (md5/sha256/sha512): ")
	actual := computeHash([]byte(text), algo)
	if actual == "" {
		fmt.Println("Неизвестный алгоритм")
		return
	}
	if strings.EqualFold(actual, expected) {
		fmt.Println("Хэши совпадают — данные целостны")
	} else {
		fmt.Println("Хэши НЕ совпадают")
		fmt.Println("ожидалось:", expected)
		fmt.Println("получено: ", actual)
	}
}
func encryptAES(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}
func decryptAES(hexStr string, key []byte) ([]byte, error) {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("шифртекст слишком короткий")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
func doAES() {
	text := readLine("Строка: ")
	keyStr := readLine("Ключ (16, 24 или 32 символа): ")
	key := []byte(keyStr)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		fmt.Println("Ключ должен быть 16, 24 или 32 байта")
		return
	}
	enc, err := encryptAES([]byte(text), key)
	if err != nil {
		fmt.Println("Ошибка шифрования:", err)
		return
	}
	fmt.Println("Зашифровано:", enc)
	decKey := readLine("Тот же ключ для расшифровки: ")
	dec, err := decryptAES(enc, []byte(decKey))
	if err != nil {
		fmt.Println("Ошибка расшифровки:", err)
		return
	}
	fmt.Println("Расшифровано:", string(dec))
}
func generateRSAKeys() {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Ошибка генерации:", err)
		return
	}
	privFile, _ := os.Create("private.pem")
	pem.Encode(privFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})
	privFile.Close()
	pubBytes, _ := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	pubFile, _ := os.Create("public.pem")
	pem.Encode(pubFile, &pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
	pubFile.Close()
	fmt.Println("Ключи сохранены: private.pem и public.pem")
}
func loadPrivateKey() (*rsa.PrivateKey, error) {
	data, err := os.ReadFile("private.pem")
	if err != nil {
		return nil, fmt.Errorf("нет private.pem — сначала сгенерируй ключи")
	}
	block, _ := pem.Decode(data)
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
func loadPublicKey() (*rsa.PublicKey, error) {
	data, err := os.ReadFile("public.pem")
	if err != nil {
		return nil, fmt.Errorf("нет public.pem — сначала сгенерируй ключи")
	}
	block, _ := pem.Decode(data)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}
func doSign() {
	msg := readLine("Сообщение: ")
	priv, err := loadPrivateKey()
	if err != nil {
		fmt.Println(err)
		return
	}
	hash := sha256.Sum256([]byte(msg))
	sig, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hash[:])
	if err != nil {
		fmt.Println("Ошибка подписи:", err)
		return
	}
	fmt.Printf("Подпись (hex): %x
", sig)
}
func doVerifySignature() {
	msg := readLine("Сообщение: ")
	sigHex := readLine("Подпись (hex): ")
	sig, err := hex.DecodeString(strings.TrimSpace(sigHex))
	if err != nil {
		fmt.Println("Подпись должна быть hex-строкой")
		return
	}
	pub, err := loadPublicKey()
	if err != nil {
		fmt.Println(err)
		return
	}
	hash := sha256.Sum256([]byte(msg))
	err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash[:], sig)
	if err != nil {
		fmt.Println("Подпись НЕВЕРНА")
		return
	}
	fmt.Println("Подпись верна")
}
func main() {
	for {
		fmt.Println()
		fmt.Println("Меню:")
		fmt.Println("1 - Вычислить хэш")
		fmt.Println("2 - Проверить целостность")
		fmt.Println("3 - AES шифрование")
		fmt.Println("4 - RSA: сгенерировать ключи")
		fmt.Println("5 - RSA: подписать")
		fmt.Println("6 - RSA: проверить подпись")
		fmt.Println("0 - Выход")
		choice := readLine("Выбор: ")
		switch choice {
		case "1":
			doHash()
		case "2":
			doVerify()
		case "3":
			doAES()
		case "4":
			generateRSAKeys()
		case "5":
			doSign()
		case "6":
			doVerifySignature()
		case "0":
			return
		}
	}
}
```

### `users-api/tlslab/cert/main.go`
```go
package main
import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)
func saveCert(certFile, keyFile string, der []byte, key *rsa.PrivateKey) {
	certOut, _ := os.Create(certFile)
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: der})
	certOut.Close()
	keyOut, _ := os.Create(keyFile)
	pem.Encode(keyOut, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	keyOut.Close()
}
func main() {
	caKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	caTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Lab CA"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}
	caDer, _ := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	saveCert("ca.crt", "ca.key", caDer, caKey)
	caCert, _ := x509.ParseCertificate(caDer)
	fmt.Println("ca.crt / ca.key")
	serverKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	serverTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: "localhost"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{"localhost"},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	serverDer, _ := x509.CreateCertificate(rand.Reader, serverTemplate, caCert, &serverKey.PublicKey, caKey)
	saveCert("server.crt", "server.key", serverDer, serverKey)
	fmt.Println("server.crt / server.key")
	clientKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	clientTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject:      pkix.Name{CommonName: "lab-client"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	clientDer, _ := x509.CreateCertificate(rand.Reader, clientTemplate, caCert, &clientKey.PublicKey, caKey)
	saveCert("client.crt", "client.key", clientDer, clientKey)
	fmt.Println("client.crt / client.key")
}
```

### `users-api/tlslab/client/main.go`
```go
package main
import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)
func main() {
	caPEM, err := os.ReadFile("../../tlslab/ca.crt")
	if err != nil {
		fmt.Println("Ошибка чтения ca.crt:", err)
		return
	}
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caPEM)
	clientCert, err := tls.LoadX509KeyPair("../../tlslab/client.crt", "../../tlslab/client.key")
	if err != nil {
		fmt.Println("Ошибка загрузки клиентского сертификата:", err)
		return
	}
	config := &tls.Config{
		RootCAs:      caPool,
		ServerName:   "localhost",
		Certificates: []tls.Certificate{clientCert},
	}
	conn, err := tls.Dial("tcp", "localhost:8443", config)
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Подключено по TLS, сертификат сервера проверен через ca.crt")
	fmt.Println("Клиентский сертификат предъявлен серверу")
	fmt.Print("Введите сообщение: ")
	reader := bufio.NewReader(os.Stdin)
	message, _ := reader.ReadString('
')
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Ошибка отправки:", err)
		return
	}
	response, err := bufio.NewReader(conn).ReadString('
')
	if err != nil {
		fmt.Println("Ошибка получения ответа:", err)
		return
	}
	fmt.Print("Ответ от сервера: ", response)
}
```

### `users-api/tlslab/server/main.go`
```go
package main
import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)
func main() {
	cert, err := tls.LoadX509KeyPair("../../tlslab/server.crt", "../../tlslab/server.key")
	if err != nil {
		fmt.Println("Ошибка загрузки сертификата:", err)
		return
	}
	caPEM, err := os.ReadFile("../../tlslab/ca.crt")
	if err != nil {
		fmt.Println("Ошибка чтения ca.crt:", err)
		return
	}
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caPEM)
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caPool,
	}
	listener, err := tls.Listen("tcp", ":8443", config)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
	fmt.Println("TLS TCP сервер запущен на :8443")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					fmt.Println("Ошибка принятия:", err)
					continue
				}
			}
			wg.Add(1)
			go func(c net.Conn) {
				defer wg.Done()
				handleConnection(c)
			}(conn)
		}
	}()
	<-stop
	fmt.Println("
Получен сигнал остановки...")
	cancel()
	listener.Close()
	wg.Wait()
	fmt.Println("Сервер остановлен.")
}
func handleConnection(conn net.Conn) {
	defer conn.Close()
	if tlsConn, ok := conn.(*tls.Conn); ok {
		if err := tlsConn.Handshake(); err != nil {
			fmt.Println("Рукопожатие не удалось:", err)
			return
		}
		state := tlsConn.ConnectionState()
		if len(state.PeerCertificates) > 0 {
			fmt.Println("Клиент:", state.PeerCertificates[0].Subject.CommonName)
		}
	}
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('
')
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
		return
	}
	fmt.Printf("Получено сообщение: %s", message)
	_, err = conn.Write([]byte("Сообщение получено по TLS!
"))
	if err != nil {
		fmt.Println("Ошибка отправки:", err)
	}
}
```

### `users-api/auth.go`
```go
package main
import (
	"encoding/json"
	"net/http"
	"time"
	"github.com/golang-jwt/jwt/v5"
)
var jwtKey = []byte("my-secret-key-for-lab")
type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
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
	role := "user"
	if id == 4 {
		role = "admin"
	}
	claims := Claims{
		UserID: id,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "jwt", "Не удалось создать токен")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
}
func checkJWT(r *http.Request) (*Claims, bool) {
	h := r.Header.Get("Authorization")
	if len(h) < 8 || h[:7] != "Bearer " {
		return nil, false
	}
	tokenStr := h[7:]
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, false
	}
	return claims, true
}
func requireAuth(w http.ResponseWriter, r *http.Request) bool {
	if _, ok := checkJWT(r); !ok {
		sendError(w, http.StatusUnauthorized, "auth", "Требуется авторизация")
		return false
	}
	return true
}
func requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	claims, ok := checkJWT(r)
	if !ok {
		sendError(w, http.StatusUnauthorized, "auth", "Требуется авторизация")
		return false
	}
	if claims.Role != "admin" {
		sendError(w, http.StatusForbidden, "forbidden", "Только для администратора")
		return false
	}
	return true
}
```

### `users-api/csrf.go`
```go
package main
import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"
)
// ===== CSRF =====
func csrfMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			token := generateRandomString(16)
			http.SetCookie(w, &http.Cookie{
				Name:     "csrf_token",
				Value:    token,
				Path:     "/",
				HttpOnly: false,
			})
			next(w, r)
			return
		}
		cookie, err := r.Cookie("csrf_token")
		header := r.Header.Get("X-CSRF-Token")
		if err != nil || header == "" || cookie.Value != header {
			sendError(w, http.StatusForbidden, "csrf", "CSRF-токен неверен")
			return
		}
		next(w, r)
	}
}
func generateRandomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}
// ===== СЕССИИ =====
type Session struct {
	UserID    int
	ExpiresAt time.Time
}
var (
	sessionsMu sync.Mutex
	sessions   = map[string]Session{}
)
func createSession(userID int) string {
	id := generateRandomString(16)
	sessionsMu.Lock()
	sessions[id] = Session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
	sessionsMu.Unlock()
	return id
}
func getSession(id string) (int, bool) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	s, ok := sessions[id]
	if !ok || time.Now().After(s.ExpiresAt) {
		delete(sessions, id)
		return 0, false
	}
	return s.UserID, true
}
```

### `users-api/db.go`
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

### `users-api/errors.go`
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

### `users-api/go.mod`
```go
module users-api
go 1.21
require (
    github.com/golang-jwt/jwt/v5 v5.3.1
    github.com/lib/pq v1.10.9
)
```

### `users-api/go.sum`
```text
github.com/golang-jwt/jwt/v5 v5.3.1 h1:kYf81DTWFe7t+1VvL7eS+jKFVWaUnK9cB1qbwn63YCY=
github.com/golang-jwt/jwt/v5 v5.3.1/go.mod h1:fxCRLWMO43lRc8nhHWY6LGqRcf+1gQWArsqaEUEa5bE=
github.com/lib/pq v1.10.9 h1:YXG7RB+JIjhP29X+OtkiDnYaXQwpS4JEWq7dtCCRUEw=
github.com/lib/pq v1.10.9/go.mod h1:AlVN5x4E4T544tWzH6hKfbfQvm3HdbOxrmggDNAPY9o=
```

### `users-api/handlers_test.go`
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

### `users-api/handlers.go`
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

### `users-api/main.go`
```go
package main
import (
	"log"
	"net/http"
	"strconv"
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
		if r.Method == http.MethodDelete {
			if !requireAdmin(w, r) {
				return
			}
		} else {
			if !requireAuth(w, r) {
				return
			}
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
	http.HandleFunc("/session/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			sendError(w, http.StatusMethodNotAllowed, "method", "Метод не разрешён")
			return
		}
		email := r.FormValue("email")
		var id int
		err := db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&id)
		if err != nil {
			sendError(w, http.StatusUnauthorized, "auth", "Не найден")
			return
		}
		sid := createSession(id)
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sid,
			Path:     "/",
			HttpOnly: true,
		})
		w.Write([]byte("Сессия создана"))
	})
	http.HandleFunc("/session/whoami", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			sendError(w, http.StatusUnauthorized, "auth", "Нет сессии")
			return
		}
		userID, ok := getSession(cookie.Value)
		if !ok {
			sendError(w, http.StatusUnauthorized, "auth", "Сессия истекла")
			return
		}
		w.Write([]byte("Ваш user_id: " + strconv.Itoa(userID)))
	})
	http.HandleFunc("/csrf-demo", csrfMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("CSRF-проверка пройдена"))
	}))
	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```
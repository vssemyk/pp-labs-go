# Лабораторная №7 - Сетевое программирование (TCP, HTTP и WebSocket) в Go

## Выполненные задания

1. Создан TCP-сервер, прослушивающий порт `:8080`, принимающий входящие соединения, считывающий сообщения от клиентов и отправляющий подтверждение о получении.
2. Реализован TCP-клиент, подключающийся к серверу, отправляющий введённое пользователем сообщение, принимающий ответ и закрывающий соединение.
3. Добавлена асинхронная обработка клиентских соединений с помощью горутин и реализован механизм корректного завершения работы (graceful shutdown) через `context`, `sync.WaitGroup` и перехват сигналов `SIGINT`/`SIGTERM`.
4. Создан базовый HTTP-сервер на порту `:8081` с обработкой маршрутов `GET /hello` и `POST /data` (декодирование JSON-данных с валидацией).
5. Добавлена маршрутизация и middleware `loggingMiddleware` для логирования HTTP-метода, пути URL и времени выполнения каждого запроса.
6. Реализован серверный WebSocket-чат на базе пакета `golang.org/x/net/websocket` с поддержкой нескольких клиентов, широковещательной рассылкой сообщений и веб-интерфейсом `chat.html`.

## Требования

- Go 1.20+
- Внешний пакет `golang.org/x/net` (`v0.61.0`) для работы с WebSocket

## Установка и запуск

Сохраните файлы согласно структуре проекта в директории `lb7`.

1. **Запуск TCP-сервера:**
Перейдите в каталог `tcpserver` и выполните команду:
```bash
go run main.go
```
Сервер выведет: `TCP сервер запущен на :8080`.

2. **Запуск TCP-клиента (в отдельном терминале):**
Перейдите в каталог `tcpclient` и выполните команду:
```bash
go run main.go
```
После запуска введите сообщение и нажмите Enter — клиент отправит его на сервер и выведет ответ.

3. **Запуск HTTP и WebSocket-сервера:**
Перейдите в каталог `httpserver`, убедитесь, что зависимости загружены (`go mod tidy`), и выполните команду:
```bash
go run main.go
```
Сервер выведет:
`HTTP сервер запущен на :8081`
`WebSocket чат: ws://localhost:8081/chat`

Для проверки чата откройте в браузере `http://localhost:8081/chat.html` в нескольких вкладках.

## Структура проекта

| Файл | Назначение |
|---|---|
| `tcpserver/main.go` | Многопоточный TCP-сервер с поддержкой graceful shutdown |
| `tcpclient/main.go` | Консольный TCP-клиент для отправки сообщений серверу |
| `httpserver/main.go` | HTTP-сервер с маршрутами `/hello`, `/data`, логирующим middleware и WebSocket-чатом `/chat` |
| `httpserver/chat.html` | Веб-клиент (HTML-страница) для тестирования WebSocket-чата |
| `httpserver/go.mod` | Файл описания модуля Go для HTTP-сервера |
| `httpserver/go.sum` | Контрольные суммы зависимостей (`golang.org/x/net`) |

## Задание 1 и 3 - TCP-сервер и асинхронная обработка с graceful shutdown
Сервер запускается на порту `:8080` через `net.Listen("tcp", ":8080")`. Каждое новое подключение обрабатывается в отдельной горутине `go func(c net.Conn)`. При получении системного сигнала остановки (`Ctrl+C` / `SIGINT` или `SIGTERM`) сервер отменяет контекст, закрывает слушатель `listener.Close()` и ожидает завершения всех активных соединений с помощью `wg.Wait()`.

### Пример
Подключение клиента, отправка сообщения `Привет, сервер!` и последующая остановка сервера по `Ctrl+C`.

### Ответ
```
TCP сервер запущен на :8080
Получено сообщение: Привет, сервер!
^C
Получен сигнал остановки. Завершаем работу...
Сервер остановлен.
```

## Задание 2 - TCP-клиент
Клиент устанавливает соединение с `localhost:8080` с помощью `net.Dial`, считывает строку из стандартного ввода, отправляет её по сети и выводит ответ сервера, после чего закрывает соединение.

### Пример
`Введите сообщение: Привет, сервер!`

### Ответ
`Ответ от сервера: Сообщение получено!`

### Тест - Сервер недоступен
Попытка запуска TCP-клиента при выключенном TCP-сервере.

### Ответ
`Ошибка подключения: dial tcp [::1]:8080: connectex: No connection could be made because the target machine actively refused it.`

## Задание 4 и 5 - HTTP-сервер, маршрутизация и middleware
HTTP-сервер прослушивает порт `:8081`. Все HTTP-обработчики обёрнуты в функцию `loggingMiddleware`, которая замеряет время выполнения запроса и записывает в журнал метод, путь и длительность обработки.
- `GET /hello` возвращает текстовое приветствие.
- `POST /data` считывает тело запроса в формате JSON в `map[string]interface{}` и выводит полученные данные в консоль сервера.

### Пример
Отправка запросов `GET http://localhost:8081/hello` и `POST http://localhost:8081/data` с телом `{"name": "Иван", "age": 25}`.

### Ответ
В ответе клиенту:
`Привет, мир!`
`Данные получены`

В консоли сервера:
```
2026/11/25 16:30:10 GET /hello 512.4µs
Получены данные: map[age:25 name:Иван]
2026/11/25 16:30:15 POST /data 1.24ms
```

### Тест - Неверный JSON
Отправка `POST /data` с некорректным телом запроса (например, `{invalid_json}`).

### Ответ
`Неверный JSON` (HTTP статус `400 Bad Request`)

## Задание 6 - WebSocket-чат
Реализован чат на основе `golang.org/x/net/websocket`. Обработчик `chatHandler` регистрирует каждого нового клиента в карте `clients`, принимает от него сообщения в цикле и отправляет их в канал `broadcast`. Фоновая горутина `handleMessages` читает сообщения из канала `broadcast` и рассылает их всем подключённым клиентам. Для работы через браузер сервер отдаёт страницу `chat.html` по маршруту `/chat.html`.

### Пример
Подключение двух клиентов через браузер (`http://localhost:8081/chat.html`) и отправка сообщения `Всем привет!`.

### Ответ
Сообщение `Всем привет!` мгновенно отображается у всех подключённых участников чата.

## Исходный код

### `tcpserver/main.go`
```go
package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
	fmt.Println("TCP сервер запущен на :8080")

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
	fmt.Println("\nПолучен сигнал остановки. Завершаем работу...")
	cancel()
	listener.Close()
	wg.Wait()
	fmt.Println("Сервер остановлен.")
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
		return
	}

	fmt.Printf("Получено сообщение: %s", message)

	_, err = conn.Write([]byte("Сообщение получено!\n"))
	if err != nil {
		fmt.Println("Ошибка отправки:", err)
	}
}
```

### `tcpclient/main.go`
```go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		return
	}
	defer conn.Close()

	fmt.Print("Введите сообщение: ")
	reader := bufio.NewReader(os.Stdin)
	message, _ := reader.ReadString('\n')

	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Ошибка отправки:", err)
		return
	}

	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка получения ответа:", err)
		return
	}

	fmt.Print("Ответ от сервера: ", response)
}
```

### `httpserver/main.go`
```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Привет, мир!")
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}
	fmt.Println("Получены данные:", data)
	fmt.Fprintln(w, "Данные получены")
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

func chatHandler(ws *websocket.Conn) {
	clients[ws] = true
	defer func() {
		delete(clients, ws)
		ws.Close()
	}()

	for {
		var msg string
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			return
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := websocket.Message.Send(client, msg)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	go handleMessages()

	http.HandleFunc("/hello", loggingMiddleware(helloHandler))
	http.HandleFunc("/data", loggingMiddleware(dataHandler))
	http.Handle("/chat", websocket.Handler(chatHandler))

	fmt.Println("HTTP сервер запущен на :8081")
	fmt.Println("WebSocket чат: ws://localhost:8081/chat")
	http.HandleFunc("/chat.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "chat.html")
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
}
```
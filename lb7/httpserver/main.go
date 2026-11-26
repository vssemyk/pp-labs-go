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

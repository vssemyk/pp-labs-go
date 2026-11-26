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

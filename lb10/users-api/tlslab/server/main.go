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
	fmt.Println("\nПолучен сигнал остановки...")
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
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
		return
	}
	fmt.Printf("Получено сообщение: %s", message)

	_, err = conn.Write([]byte("Сообщение получено по TLS!\n"))
	if err != nil {
		fmt.Println("Ошибка отправки:", err)
	}
}

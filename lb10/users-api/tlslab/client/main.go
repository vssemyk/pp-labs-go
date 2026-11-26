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

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
	line, _ := reader.ReadString('\n')
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
	fmt.Printf("Хэш (%s): %s\n", algo, result)
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
	fmt.Printf("Подпись (hex): %x\n", sig)
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

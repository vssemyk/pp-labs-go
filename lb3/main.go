package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"lb3/mathutils"
	"lb3/stringutils"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Задание 2: факториал введённого числа
	fmt.Println("=== 1. Вычисление факториала ===")
	fmt.Print("Введите целое число: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	n, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Ошибка: введено не целое число")
	} else {
		fact, err := mathutils.Factorial(n)
		if err != nil {
			fmt.Printf("Ошибка: %v\n", err)
		} else {
			fmt.Printf("%d! = %d\n", n, fact)
		}
	}

	// Задание 3: переворот строки
	fmt.Println("\n=== 2. Переворот строки ===")
	fmt.Print("Введите строку: ")
	str, _ := reader.ReadString('\n')
	str = strings.TrimRight(str, "\r\n")
	fmt.Printf("Перевёрнутая строка: %s\n", stringutils.Reverse(str))

	// Задание 4: массив из 5 целых чисел
	fmt.Println("\n=== 3. Массив из 5 целых чисел ===")
	var arr [5]int

	for i := 0; i < len(arr); i++ {
		arr[i] = rand.Intn(100)
		fmt.Printf("arr[%d] = %d\n", i, arr[i])
	}

	fmt.Print("Массив: ")
	for _, v := range arr {
		fmt.Printf("%d ", v)
	}
	fmt.Println()

	// Задание 5: срез из массива, добавление и удаление
	fmt.Println("\n=== 4. Работа со срезом ===")
	slice := arr[:]
	fmt.Printf("Исходный срез: %v \n", slice)

	// Добавление элемента
	slice = append(slice, 999)
	fmt.Printf("После append(999): %v \n", slice)

	// Удаление элемента по индексу 2
	indexToRemove := 2
	if indexToRemove >= 0 && indexToRemove < len(slice) {
		slice = append(slice[:indexToRemove], slice[indexToRemove+1:]...)
		fmt.Printf("После удаления элемента с индексом %d: %v \n",
			indexToRemove, slice)
	}

	// Задание 6: самая длинная строка в срезе
	fmt.Println("\n=== 5. Поиск самой длинной строки ===")
	words := []string{"привет", "бронетранспортер", "привычка", "котенок", "ответ"}
	fmt.Printf("Срез строк: %v\n", words)
	longest := stringutils.LongestString(words)
	fmt.Printf("Самая длинная строка: %q (длина = %d)\n", longest, len([]rune(longest)))
}

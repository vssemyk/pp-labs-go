package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Функция среднего возраста
func averageAge(people map[string]int) float64 {
	if len(people) == 0 {
		return 0
	}
	sum := 0
	for _, age := range people {
		sum += age
	}
	return float64(sum) / float64(len(people))
}

func task1() {
	fmt.Println("\n--- Задание 1: Карта с именами и возрастами ---")
	people := map[string]int{
		"Анна":   25,
		"Борис":  30,
		"Виктор": 22,
	}

	people["Галина"] = 28

	for name, age := range people {
		fmt.Printf("%s: %d лет\n", name, age)
	}
}

// Удаление записи по имени
func task3() {
	fmt.Println("\n--- Задание 3: Удаление записи по имени ---")
	people := map[string]int{
		"Анна":   25,
		"Борис":  30,
		"Виктор": 22,
	}

	nameToDelete := "Борис"
	fmt.Printf("Удаляем: %s\n", nameToDelete)

	if _, exists := people[nameToDelete]; exists {
		delete(people, nameToDelete)
		fmt.Println("Запись удалена.")
	} else {
		fmt.Println("Запись не найдена.")
	}

	fmt.Println("Оставшиеся записи:")
	for name, age := range people {
		fmt.Printf("%s: %d лет\n", name, age)
	}
}

// Вывод строки в верхнем регистре
func task4() {
	fmt.Println("\n--- Задание 4: Строка в верхнем регистре ---")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите строку: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	fmt.Println("В верхнем регистре:", strings.ToUpper(input))
}

// Сумма нескольких чисел
func task5() {
	fmt.Println("\n--- Задание 5: Сумма чисел ---")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите числа через пробел: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	fields := strings.Fields(input)
	sum := 0
	for _, f := range fields {
		num, err := strconv.Atoi(f)
		if err != nil {
			fmt.Printf("Пропущено нечисловое значение: %s\n", f)
			continue
		}
		sum += num
	}
	fmt.Printf("Сумма: %d\n", sum)
}

// Массив в обратном порядке
func task6() {
	fmt.Println("\n--- Задание 6: Массив в обратном порядке ---")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите целые числа через пробел: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	fields := strings.Fields(input)
	arr := make([]int, 0, len(fields))
	for _, f := range fields {
		if num, err := strconv.Atoi(f); err == nil {
			arr = append(arr, num)
		}
	}

	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	fmt.Println("В обратном порядке:", arr)
}

func main() {
	task1()

	fmt.Println("\n--- Задача 2: Средний возраст ---")
	people := map[string]int{
		"Анна":   25,
		"Борис":  30,
		"Виктор": 22,
	}
	fmt.Printf("Средний возраст: %.2f\n", averageAge(people))

	task3()
	task4()
	task5()
	task6()
}

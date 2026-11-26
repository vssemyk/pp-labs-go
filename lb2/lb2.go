package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type rect struct {
	width, length float64
}

func isEven(num int) {
	if num%2 == 0 {
		fmt.Println(num, " - четное число")
	} else {
		fmt.Println(num, " - нечетное число")
	}

}

func isPositive(num int) string {
	if num > 0 {
		return "Positive"
	} else if num == 0 {
		return "Zero"
	} else {
		return "Negative"
	}
}

func stringLength(str string) int {
	return len([]rune(str))
}

func (rec rect) rectangleSquare() float64 {
	return rec.width * rec.length
}

func average(a, b int) float64 {
	return float64(a+b) / 2
}

func main() {

	reader := bufio.NewReader(os.Stdin)

	for {

		var num int = 0
		var option int
		fmt.Println("\n=====Выберите задание:===== \n 1. Четность\n 2. Положительность\n 3. Цикл for\n 4. Длина строки\n 5. Площадь прямоугольника\n 6. Среднее из двух\n 0. Закрыть программу")
		fmt.Scanf("%d\n", &option)

		switch option {
		case 1:

			fmt.Println("*****Определение четности числа*****")
			fmt.Print("Введите число: ")
			fmt.Scanf("%d\n", &num)

			isEven(num)

		case 2:
			fmt.Println("*****Определение положительности числа*****")
			fmt.Print("Введите число: ")
			fmt.Scanf("%d\n", &num)

			fmt.Println(isPositive(num))

		case 3:
			fmt.Println("*****Цикл for*****")

			for i := 1; i <= 10; i++ {
				fmt.Println(i)
			}

		case 4:

			fmt.Println("*****Определение длины строки*****")
			_, _ = reader.ReadString('\n')
			fmt.Print("Введите строку: ")
			text, _ := reader.ReadString('\n')
			text = strings.TrimRight(text, "\r\n")

			fmt.Println("Длина строки:", stringLength(text))

		case 5:

			var rec rect

			fmt.Println("*****Площадь прямоугольника*****")
			fmt.Print("Введите длину: ")
			fmt.Scanf("%f\n", &rec.length)
			fmt.Print("Введите ширину: ")
			fmt.Scanf("%f\n", &rec.width)

			fmt.Println("Ширина = ", rec.width, ", длина = ", rec.length)

			fmt.Println("Площадь прямоугольника: ", rec.rectangleSquare())

		case 6:
			var num1, num2 int
			fmt.Println("*****Среднее значение двух чисел*****")
			fmt.Print("Введите первое число: ")
			fmt.Scanf("%d\n", &num1)
			fmt.Print("Введите второе число: ")
			fmt.Scanf("%d\n", &num2)

			fmt.Println("Среднее: ", average(num1, num2))

		case 0:
			fmt.Println("Выход")
			return

		default:
			fmt.Println("\nНекорректный пункт")
		}
	}

}

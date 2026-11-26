package main

import (
	"fmt"
	"time"
)

func sum(a, b float64) float64 {
	return a + b
}

func diff(a, b float64) float64 {
	return a - b
}

func average(a, b, c float64) float64 {
	return (a + b + c) / 3
}

func main() {

	fmt.Println("=== Задание 1 ===")
	now := time.Now()
	fmt.Println("Дата и время:", now.Format("02.01.2006 15:04:05"))

	fmt.Println("\n=== Задание 2 ===")
	var i int = 10
	var f float64 = 3.14
	var s string = "Привет"
	var b bool = true
	fmt.Println("int:", i)
	fmt.Println("float64:", f)
	fmt.Println("string:", s)
	fmt.Println("bool:", b)

	fmt.Println("\n=== Задание 3 ===")
	x := 5
	y := 2.5
	name := "Go"
	ok := false
	fmt.Println("x =", x)
	fmt.Println("y =", y)
	fmt.Println("name =", name)
	fmt.Println("ok =", ok)

	fmt.Println("\n=== Задание 4 ===")
	a := 17
	c := 5
	fmt.Println("Сложение:", a+c)
	fmt.Println("Вычитание:", a-c)
	fmt.Println("Умножение:", a*c)
	fmt.Println("Деление:", a/c)
	fmt.Println("Остаток:", a%c)

	fmt.Println("\n=== Задание 5 ===")
	fmt.Println("Сумма:", sum(10.5, 4.2))
	fmt.Println("Разность:", diff(10.5, 4.2))

	fmt.Println("\n=== Задание 6 ===")
	fmt.Println("Среднее:", average(4, 7, 10))
}

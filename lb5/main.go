package main

import (
	"fmt"
	"math"
)

//структура Person и ее методы

type Person struct {
	name string
	age  int
}

func (p Person) print() {
	fmt.Printf("Имя: %s, Возраст: %d\n", p.name, p.age)
}

func (p *Person) birthday() {
	p.age++
}

//структура Circle

type Circle struct {
	radius float64
}

type Rectangle struct {
	length float64
	width  float64
}

// интерфейс Shape с методом Area()
type Shape interface {
	Area() float64
}

func (c Circle) Area() float64 {
	return c.radius * c.radius * math.Pi
}

func (r Rectangle) Area() float64 {
	return r.length * r.width
}

// функция, принимает срез интерфейсов Shape и выводит площадь каждого объекта.
func printAreas(shapes []Shape) {
	for i, shape := range shapes {
		fmt.Printf("Фигура %d: площадь = %.2f\n", i+1, shape.Area())
	}
}

// интерфейс Stringer и реализация его для структуры Book

type Stringer interface {
	String() string
}

type Book struct {
	author string
	title  string
	genre  string
	year   int
}

func (b Book) String() string {
	return fmt.Sprintf("Книга: \"%s\", автор: %s, год: %d, жанр: %s", b.title, b.author, b.year, b.genre)
}

func printStringer(s Stringer) {
	fmt.Println(s.String())
}

func main() {

	//Person
	fmt.Println("\n---Задание 1 и 2---")
	person := Person{"Влада", 26}
	fmt.Print("Текущая запись: ")
	person.print()
	person.birthday()
	fmt.Print("После дня рождения: ")
	person.print()

	//Circle
	fmt.Println("\n---Задание 3---")
	circle := Circle{5}
	fmt.Println("Площадь круга: ", circle.Area())

	//Shape
	fmt.Println("\n---Задание 4 и 5---")
	rect := Rectangle{4, 6}
	shapes := []Shape{rect, circle}
	printAreas(shapes)

	//Stringer
	fmt.Println("\n---Задание 6---")
	book := Book{"Стивен Кинг", "Зеленая миля", "драма", 1996}
	printStringer(book)

}

# Лабораторная №5 - Структуры, методы и интерфейсы в Go

## Выполненные задания

1. Создана структура `Person` с полями `name` и `age`, реализован метод `print()` для вывода информации о человеке.
2. Реализован метод `birthday()` с получателем-указателем (`*Person`), увеличивающий возраст человека на 1 год.
3. Создана структура `Circle` с полем `radius` и метод `Area()` для вычисления площади круга.
4. Создан интерфейс `Shape` с методом `Area()`, который реализован для структур `Rectangle` и `Circle`.
5. Реализована функция `printAreas(shapes []Shape)`, принимающая срез интерфейсов `Shape` и выводящая площадь каждой фигуры.
6. Создан интерфейс `Stringer` с методом `String()` и реализован для структуры `Book`, хранящей информацию об авторе, названии, жанре и годе издания книги.

## Требования

- Go 1.20+

## Установка и запуск

Сохраните файл под именем `main.go`, перейдите в каталог с ним и выполните команду:

```bash
go run main.go
```

После запуска программа автоматически выполнит все задания по очереди и выведет результаты работы со структурами, методами и интерфейсами в консоль.

## Структура проекта

| Файл | Назначение |
|---|---|
| `main.go` | Реализация структур `Person`, `Circle`, `Rectangle`, `Book`, интерфейсов `Shape`, `Stringer` и демонстрация их работы |

## Задание 1 и 2 - Структура Person и методы print и birthday
Создана структура `Person` с полями `name string` и `age int`. Метод `print()` выводит текущие имя и возраст, а метод `birthday()` принимает указатель `*Person` и увеличивает значение поля `age` на единицу. В `main` создается экземпляр `Person{"Влада", 26}`, выводится исходное состояние, вызывается `birthday()` и выводится обновленный возраст.

### Пример
Создание экземпляра `Person{"Влада", 26}` и вызов методов `print()` и `birthday()`.

### Ответ
```
---Задание 1 и 2---
Текущая запись: Имя: Влада, Возраст: 26
После дня рождения: Имя: Влада, Возраст: 27
```

## Задание 3 - Структура Circle и вычисление площади круга
Создана структура `Circle` с полем `radius float64` и метод `Area() float64`, вычисляющий площадь круга по формуле $\pi \cdot r^2$ с использованием константы `math.Pi`.

### Пример
Создание круга с радиусом `Circle{5}` и вызов метода `circle.Area()`.

### Ответ
```
---Задание 3---
Площадь круга:  78.53981633974483
```

## Задание 4 и 5 - Интерфейс Shape и полиморфный вывод площадей
Определен интерфейс `Shape`, требующий наличия метода `Area() float64`. Помимо `Circle`, этот интерфейс реализован для структуры `Rectangle` (поля `length` и `width`). Функция `printAreas(shapes []Shape)` перебирает срез фигур и выводит площадь каждой с округлением до двух знаков после запятой.

### Пример
Передача среза `[]Shape{Rectangle{4, 6}, Circle{5}}` в функцию `printAreas`.

### Ответ
```
---Задание 4 и 5---
Фигура 1: площадь = 24.00
Фигура 2: площадь = 78.54
```

## Задание 6 - Интерфейс Stringer и структура Book
Создан интерфейс `Stringer` с методом `String() string` и структура `Book` с полями `author`, `title`, `genre`, `year`. Реализован метод `String()` для `Book`, формирующий читаемую строку с описанием книги, а функция `printStringer(s Stringer)` выводит эту строку в консоль.

### Пример
Создание объекта `Book{"Стивен Кинг", "Зеленая миля", "драма", 1996}` и передача в `printStringer`.

### Ответ
```
---Задание 6---
Книга: "Зеленая миля", автор: Стивен Кинг, год: 1996, жанр: драма
```

## Исходный код

```go
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
```
# Лабораторная №6 - Многопоточность, горутины и каналы в Go

## Выполненные задания

1. Создана программа с тремя параллельно выполняющимися горутинами: расчёт факториала, генерация случайных чисел и вычисление суммы числового ряда с имитацией задержек через `time.Sleep()`.
2. Реализована передача первых 10 чисел Фибоначчи через канал с последующим закрытием канала функцией `close()` и чтением через цикл `for range`.
3. Продемонстрирована работа конструкции `select` для одновременного чтения данных из двух каналов (генератор случайных чисел и определитель чётности).
4. Реализована синхронизация доступа к общей переменной-счётчику с помощью мьютекса (`sync.Mutex`) и показано возникновение гонки данных (data race) при его отключении.
5. Разработан многопоточный калькулятор на базе горутин и каналов, поддерживающий операции `+`, `-`, `*`, `/` и обработку деления на ноль.
6. Реализован пул воркеров (Worker Pool) с задаваемым пользователем количеством обработчиков для параллельного реверсирования строк через буферизованные каналы.

## Требования

- Go 1.20+

## Установка и запуск

Сохраните файл под именем `main.go`, перейдите в каталог с ним и выполните команду:

```bash
go run main.go
```

После запуска программа последовательно выполнит задания 1–5 в автоматическом режиме, а в задании 6 запросит ввод количества воркеров с клавиатуры и выведет результаты обработки строк.

## Структура проекта

| Файл | Назначение |
|---|---|
| `main.go` | Демонстрация работы горутин, каналов, `select`, мьютексов, многопоточного калькулятора и пула воркеров |

## Задание 1 - Создание и запуск горутин
Запускает в отдельных горутинах три независимые функции: `task1_factorial(5)`, `task1_random(5)` и `task1_sum(10)`. За счёт использования `time.Sleep()` функции выполняются конкурентно и выводят результаты по мере готовности, а основная горутина ожидает их завершения 3 секунды.

### Пример
Автоматический запуск трёх горутин из функции `runTask1()`.

### Ответ
```
--- Задание 1 ---
Случайное число: 81
Факториал 5 = 120
Случайное число: 87
Случайное число: 47
Сумма ряда = 2.9289682539682538
Случайное число: 59
Случайное число: 81
Готово
```

## Задание 2 - Использование каналов для передачи данных
Горутина `fibonacci(10, ch)` вычисляет первые 10 чисел Фибоначчи и отправляет их в небуферизованный канал `ch`, после чего закрывает его вызовом `close(ch)`. Основная горутина считывает значения в цикле `for num := range ch`. Вызов `close(ch)` сигнализирует получателю о том, что новых данных больше не будет, благодаря чему цикл `range` корректно завершается без взаимной блокировки (deadlock).

### Пример
Генерация последовательности из `10` чисел Фибоначчи.

### Ответ
```
--- Задание 2 ---
Получено: 0
Получено: 1
Получено: 1
Получено: 2
Получено: 3
Получено: 5
Получено: 8
Получено: 13
Получено: 21
Получено: 34
Канал закрыт
```

## Задание 3 - Применение select для управления каналами
Запускаются две фоновые горутины: первая генерирует случайные числа и отправляет их в канал `numbersCh`, вторая считывает числа из `numbersCh` и отправляет строку о чётности в `parityCh`. В цикле выполняется 10 итераций конструкции `select`, которая неблокирующим образом выбирает готовый к чтению канал и выводит полученное сообщение.

### Пример
10 итераций опроса каналов `numbersCh` и `parityCh` через `select`.

### Ответ
```
--- Задание 3 ---
Сгенерировано: 40
Сгенерировано: 56
Результат: Число 94 — чётное
Сгенерировано: 87
Результат: Число 25 — нечётное
Сгенерировано: 28
Сгенерировано: 32
Результат: Число 37 — нечётное
Сгенерировано: 91
Результат: Число 60 — чётное
```

## Задание 4 - Синхронизация с помощью мьютексов
Пять горутин по 10 000 раз инкрементируют общую переменную `counter` (ожидаемое итоговое значение — `50000`). В первом проходе мьютекс выключен (`useMutex = false`), из-за чего возникает состояние гонки (data race) и часть инкрементов теряется. Во втором проходе доступ защищён `mu.Lock()` и `mu.Unlock()`, что гарантирует точный результат.

### Пример
Запуск 5 горутин по 10 000 итераций без `sync.Mutex` и с использованием `sync.Mutex`.

### Ответ
```
--- Задание 4 ---
Без мьютекса counter = 41290 (ожидалось 50000)
С мьютексом counter = 50000 (ожидалось 50000)
```

## Задание 5 - Разработка многопоточного калькулятора
Серверная горутина `calculatorServer` принимает структуры `Request` (операнды `A`, `B`, операция `Op` и канал ответа `Result`) из канала `requests`, вычисляет результат и отправляет его обратно клиенту. Поддерживаются операции `+`, `-`, `*`, `/` с проверкой деления на ноль.

### Пример
Отправка запросов на вычисление: `10 + 5`, `10 - 5`, `10 * 5`, `10 / 5`.

### Ответ
```
--- Задание 5 ---
10.0 + 5.0 = 15.00
10.0 - 5.0 = 5.00
10.0 * 5.0 = 50.00
10.0 / 5.0 = 2.00
```

### Тест - Деление на ноль
`Request{A: 10, B: 0, Op: "/", Result: resultCh}`

### Ответ
```
Ошибка: деление на ноль
10.0 / 0.0 = 0.00
```

## Задание 6 - Создание пула воркеров
Создаётся пул из заданного пользователем количества горутин-воркеров (`worker`), которые читают строки из канала `jobs`, переворачивают их с учётом рун (`[]rune`) и отправляют результат в канал `results`. После завершения работы всех воркеров (`wg.Wait()`) канал `results` закрывается, и программа выводит итоги.

### Пример
`Введите количество воркеров: 3`

### Ответ
```
Результаты:
Воркер 3: зачем -> мечаз
Воркер 1: почему -> умечоп
Воркер 2: как -> как
Воркер 3: планета -> атеналп
Воркер 1: сатурн -> нрутас
Воркер 2: венера -> аренев
Воркер 3: марс -> срам
```

## Исходный код

```go
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ЗАДАНИЕ 1
func task1_factorial(n int) {
	result := 1
	for i := 1; i <= n; i++ {
		result *= i
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("Факториал", n, "=", result)
}

func task1_random(count int) {
	for i := 0; i < count; i++ {
		fmt.Println("Случайное число:", rand.Intn(100))
		time.Sleep(150 * time.Millisecond)
	}
}

func task1_sum(n int) {
	sum := 0.0
	for i := 1; i <= n; i++ {
		sum += 1.0 / float64(i)
		time.Sleep(80 * time.Millisecond)
	}
	fmt.Println("Сумма ряда =", sum)
}

func runTask1() {
	fmt.Println("--- Задание 1 ---")
	go task1_factorial(5)
	go task1_random(5)
	go task1_sum(10)
	time.Sleep(3 * time.Second)
	fmt.Println("Готово")
}

// ЗАДАНИЕ 2
func fibonacci(n int, ch chan int) {
	a, b := 0, 1
	for i := 0; i < n; i++ {
		ch <- a
		a, b = b, a+b
		time.Sleep(200 * time.Millisecond)
	}
	close(ch)
}

func runTask2() {
	fmt.Println("--- Задание 2 ---")
	ch := make(chan int)
	go fibonacci(10, ch)
	for num := range ch {
		fmt.Println("Получено:", num)
	}
	fmt.Println("Канал закрыт")
}

// ЗАДАНИЕ 3
func runTask3() {
	fmt.Println("--- Задание 3 ---")
	numbersCh := make(chan int)
	parityCh := make(chan string)

	go func() {
		for {
			n := rand.Intn(100)
			numbersCh <- n
			time.Sleep(300 * time.Millisecond)
		}
	}()

	go func() {
		for {
			n := <-numbersCh
			if n%2 == 0 {
				parityCh <- fmt.Sprintf("Число %d — чётное", n)
			} else {
				parityCh <- fmt.Sprintf("Число %d — нечётное", n)
			}
		}
	}()

	for i := 0; i < 10; i++ {
		select {
		case n := <-numbersCh:
			fmt.Println("Сгенерировано:", n)
		case msg := <-parityCh:
			fmt.Println("Результат:", msg)
		}
	}
}

// ЗАДАНИЕ 4
var counter int
var mu sync.Mutex

func increment(wg *sync.WaitGroup, useMutex bool) {
	defer wg.Done()
	for i := 0; i < 10000; i++ {
		if useMutex {
			mu.Lock()
		}
		counter++
		if useMutex {
			mu.Unlock()
		}
	}
}

func runTask4() {
	fmt.Println("--- Задание 4 ---")

	counter = 0
	var wg1 sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg1.Add(1)
		go increment(&wg1, false)
	}
	wg1.Wait()
	fmt.Println("Без мьютекса counter =", counter, "(ожидалось 50000)")

	counter = 0
	var wg2 sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg2.Add(1)
		go increment(&wg2, true)
	}
	wg2.Wait()
	fmt.Println("С мьютексом counter =", counter, "(ожидалось 50000)")
}

// ЗАДАНИЕ 5
type Request struct {
	A, B   float64
	Op     string
	Result chan float64
}

func calculatorServer(requests chan Request) {
	for req := range requests {
		var result float64
		switch req.Op {
		case "+":
			result = req.A + req.B
		case "-":
			result = req.A - req.B
		case "*":
			result = req.A * req.B
		case "/":
			if req.B != 0 {
				result = req.A / req.B
			} else {
				fmt.Println("Ошибка: деление на ноль")
			}
		}
		req.Result <- result
	}
}

func runTask5() {
	fmt.Println("--- Задание 5 ---")
	requests := make(chan Request)
	go calculatorServer(requests)

	ops := []struct {
		a, b float64
		op   string
	}{
		{10, 5, "+"},
		{10, 5, "-"},
		{10, 5, "*"},
		{10, 5, "/"},
	}

	for _, o := range ops {
		resultCh := make(chan float64)
		requests <- Request{A: o.a, B: o.b, Op: o.op, Result: resultCh}
		res := <-resultCh
		fmt.Printf("%.1f %s %.1f = %.2f
", o.a, o.op, o.b, res)
	}
	close(requests)
}

// ЗАДАНИЕ 6
func worker(id int, jobs <-chan string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for str := range jobs {
		runes := []rune(str)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		results <- fmt.Sprintf("Воркер %d: %s -> %s", id, str, string(runes))
	}
}

func runTask6() {
	fmt.Println("--- Задание 6 ---")
	lines := []string{"почему", "как", "зачем", "планета", "сатурн", "венера", "марс"}

	var numWorkers int
	fmt.Print("Введите количество воркеров: ")
	fmt.Scan(&numWorkers)

	jobs := make(chan string, len(lines))
	results := make(chan string, len(lines))

	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, results, &wg)
	}

	for _, line := range lines {
		jobs <- line
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println("
Результаты:")
	for res := range results {
		fmt.Println(res)
	}
}

func main() {
	runTask1()
	runTask2()
	runTask3()
	runTask4()
	runTask5()
	runTask6()
}
```
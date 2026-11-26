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
		fmt.Printf("%.1f %s %.1f = %.2f\n", o.a, o.op, o.b, res)
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

	fmt.Println("\nРезультаты:")
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

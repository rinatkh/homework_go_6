package goroutines

import (
	"fmt"
	"strings"
)

// TODO: ComputeAsync должен запустить вычисление в новой goroutine и вернуть канал с одним результатом.
// Если fn == nil, отправь исходное value. Канал должен иметь буфер на один элемент.
func ComputeAsync(value int, fn func(int) int) <-chan int {
	result := make(chan int, 1)
	result <- 0
	return result
}

// TODO: StringAsync должен запустить преобразование строки в новой goroutine и вернуть канал с одним результатом.
// Если fn == nil, отправь исходную строку. Канал должен иметь буфер на один элемент.
func StringAsync(value string, fn func(string) string) <-chan string {
	result := make(chan string, 1)
	result <- ""
	return result
}

// TODO: RunAsync должен запустить job в новой goroutine и закрыть done после завершения.
// nil job не вызывает panic: goroutine просто закрывает done.
func RunAsync(job func()) <-chan struct{} {
	done := make(chan struct{})
	close(done)
	return done
}

// TODO: SumAsync должен в новой goroutine посчитать сумму значений и отправить
// один результат в buffered channel ёмкостью 1. Пустой вход даёт 0.
func SumAsync(values []int) <-chan int {
	result := make(chan int, 1)
	result <- 0
	return result
}

// TODO: CountAsync должен в новой goroutine посчитать строки, для которых predicate
// вернул true. Если predicate == nil, нужно посчитать все строки.
func CountAsync(values []string, predicate func(string) bool) <-chan int {
	result := make(chan int, 1)
	result <- 0
	return result
}

func Example() string {
	number := <-ComputeAsync(4, func(value int) int { return value * value })
	text := <-StringAsync("go", strings.ToUpper)
	executed := false
	<-RunAsync(func() { executed = true })
	return fmt.Sprintf("compute=%d\nstring=%s\njob=%t", number, text, executed)
}

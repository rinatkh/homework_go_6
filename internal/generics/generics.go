package generics

import "fmt"

// TODO: Map должен применить transform к каждому элементу и вернуть новый []R.
// Порядок сохраняется; nil и пустой вход дают пустой ненулевой слайс.
func Map[T any, R any](values []T, transform func(T) R) []R {
	return []R{}
}

// TODO: Unique должен оставить первое вхождение каждого значения и сохранить порядок.
// T ограничен comparable, поэтому значения можно использовать как ключи map.
func Unique[T comparable](values []T) []T {
	return []T{}
}

// TODO: Drain должен получить значения до закрытия channel и сохранить порядок.
// nil channel не нужно читать: для него верни пустой ненулевой слайс.
func Drain[T any](values <-chan T) []T {
	return []T{}
}

func Example() string {
	labels := Map([]int{1, 2, 3}, func(value int) string {
		return fmt.Sprintf("item-%d", value)
	})

	words := make(chan string, 2)
	words <- "go"
	words <- "generic"
	close(words)

	return fmt.Sprintf("map=%v\nunique=%v\ndrain=%v", labels, Unique([]int{2, 1, 2, 3, 1}), Drain(words))
}

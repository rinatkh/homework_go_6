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

// TODO: Filter должен вернуть новый слайс только со значениями, для которых
// keep вернул true. При keep == nil вернуть независимую копию values.
func Filter[T any](values []T, keep func(T) bool) []T {
	return []T{}
}

// TODO: IndexOf должен вернуть индекс первого target или -1, если target нет.
// T ограничен comparable, потому что значения сравниваются через ==.
func IndexOf[T comparable](values []T, target T) int {
	return -1
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

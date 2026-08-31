package waitgroup

import (
	"fmt"
	"sync/atomic"
)

// TODO: RunAll должен запустить каждую ненулевую задачу в отдельной goroutine
// и вернуться только после завершения всех задач.
func RunAll(tasks []func()) {}

// TODO: Squares должен вычислить квадраты конкурентно и вернуть их в исходном порядке.
// nil и пустой вход дают пустой ненулевой слайс.
func Squares(values []int) []int {
	return []int{}
}

// TODO: ApplyAll должен применить fn к элементам конкурентно и сохранить порядок.
// Если fn == nil, верни независимую копию values.
func ApplyAll(values []int, fn func(int) int) []int {
	return []int{}
}

// TODO: SumParts должен запустить отдельную goroutine для каждой части,
// посчитать сумму внутри части и сохранить суммы в исходном порядке частей.
func SumParts(parts [][]int) []int {
	return []int{}
}

// TODO: CountMatches должен конкурентно посчитать количество target в каждой
// группе строк и вернуть счётчики в исходном порядке групп.
func CountMatches(groups [][]string, target string) []int {
	return []int{}
}

func Example() string {
	var count atomic.Int32
	RunAll([]func(){func() { count.Add(1) }, func() { count.Add(1) }})
	return fmt.Sprintf("tasks=%d\nsquares=%v\napplied=%v", count.Load(), Squares([]int{2, 3}), ApplyAll([]int{1, 2}, func(v int) int { return v + 10 }))
}

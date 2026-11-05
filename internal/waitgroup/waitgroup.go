package waitgroup

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// TODO: RunAll должен запустить каждую ненулевую задачу в отдельной горутине
// и вернуться только после завершения всех задач.
func RunAll(tasks []func()) {
	var wg sync.WaitGroup

	for _, task := range tasks {
		if task == nil {
			continue
		}
		wg.Add(1)
		go func(t func()) {
			defer wg.Done()
			t()
		}(task)
	}
	wg.Wait()
}

// TODO: Squares должен вычислить квадраты конкурентно и вернуть их в исходном порядке.
// nil и пустой вход дают пустой ненулевой слайс.
func Squares(values []int) []int {
	result := make([]int, len(values))
	var wg sync.WaitGroup

	for index, value := range values {
		wg.Add(1)
		go func(i, v int) {
			defer wg.Done()
			result[i] = v * v
		}(index, value)
	}
	wg.Wait()
	return result
}

// TODO: ApplyAll должен применить fn к элементам конкурентно и сохранить порядок.
// Если fn == nil, верни независимую копию values.
func ApplyAll(values []int, fn func(int) int) []int {
	var wg sync.WaitGroup
	if fn == nil {
		newValues := make([]int, len(values))
		copy(newValues, values)
		return newValues
	}

	result := make([]int, len(values))

	for index, value := range values {
		wg.Add(1)
		go func(i, v int) {
			defer wg.Done()
			result[i] = fn(v)
		}(index, value)
	}
	wg.Wait()
	return result
}

// TODO: SumParts должен запустить отдельную горутину для каждой части,
// посчитать сумму внутри части и сохранить суммы в исходном порядке частей.
func SumParts(parts [][]int) []int {
	var wg sync.WaitGroup
	result := make([]int, len(parts))

	for index, part := range parts {
		wg.Add(1)
		go func(i int, p []int) {
			sum := 0
			defer wg.Done()
			for _, row := range p {
				sum = sum + row
			}
			result[i] = sum
		}(index, part)
	}
	wg.Wait()
	return result
}

// TODO: CountMatches должен конкурентно посчитать количество target в каждой
// группе строк и вернуть счётчики в исходном порядке групп.
func CountMatches(groups [][]string, target string) []int {
	var wg sync.WaitGroup
	result := make([]int, len(groups))

	for index, group := range groups {
		wg.Add(1)
		go func(i int, g []string) {
			count := 0
			defer wg.Done()
			for _, row := range g {
				if row == target {
					count++
				}
			}
			result[i] = count
		}(index, group)
	}
	wg.Wait()
	return result
}

func Example() string {
	var count atomic.Int32
	RunAll([]func(){func() { count.Add(1) }, func() { count.Add(1) }})
	return fmt.Sprintf("tasks=%d\nsquares=%v\napplied=%v", count.Load(), Squares([]int{2, 3}), ApplyAll([]int{1, 2}, func(v int) int { return v + 10 }))
}

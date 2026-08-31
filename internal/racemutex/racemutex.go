package racemutex

import (
	"fmt"
	"sync"
)

type Counter struct {
	mu    sync.Mutex
	value int
}

// TODO: Add должен безопасно прибавить delta к общему значению.
// nil receiver нужно спокойно игнорировать.
func (c *Counter) Add(delta int) {}

// TODO: Value должен безопасно вернуть текущее значение.
// nil receiver возвращает 0.
func (c *Counter) Value() int {
	return 0
}

// TODO: ParallelAdd должен запустить workers goroutines, каждая выполняет
// increments вызовов Add(1), дождаться всех и вернуть итог.
// Неположительные workers или increments дают 0.
func ParallelAdd(workers, increments int) int {
	return 0
}

// TODO: Swap должен под одной блокировкой заменить значение и вернуть прежнее.
// Для nil receiver вернуть 0 и ничего не делать.
func (c *Counter) Swap(value int) int {
	return 0
}

// TODO: ParallelDeltas должен запустить отдельную goroutine для каждого delta,
// безопасно применить все изменения к одному Counter, дождаться их и вернуть итог.
func ParallelDeltas(deltas []int) int {
	return 0
}

func Example() string {
	return fmt.Sprintf("counter=%d", ParallelAdd(8, 50))
}

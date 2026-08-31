package closerange

import "fmt"

// TODO: Generate должен в отдельной goroutine отправить числа от 0 до limit-1,
// затем закрыть канал. Неположительный limit даёт сразу закрытый канал.
func Generate(limit int) <-chan int {
	values := make(chan int)
	close(values)
	return values
}

// TODO: Collect должен прочитать канал через range до закрытия.
// nil channel трактуется как пустой вход; результат всегда ненулевой слайс.
func Collect(values <-chan int) []int {
	return []int{}
}

// TODO: FilterEven должен отправить только чётные значения в исходном порядке
// и закрыть выходной канал после завершения.
func FilterEven(values []int) <-chan int {
	result := make(chan int)
	close(result)
	return result
}

func Example() string {
	return fmt.Sprintf("generated=%v\neven=%v", Collect(Generate(5)), Collect(FilterEven([]int{1, 2, 3, 4, 5, 6})))
}

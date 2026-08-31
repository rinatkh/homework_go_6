package channels

import "fmt"

// TODO: SendPair должен вернуть buffered channel ёмкостью 2
// с first и second в таком порядке. Канал закрывать не нужно.
func SendPair(first, second int) <-chan int {
	values := make(chan int, 2)
	values <- 0
	values <- 0
	return values
}

// TODO: SumN должен получить ровно count значений и вернуть их сумму.
// count <= 0 возвращает 0 и ничего не читает из канала.
func SumN(values <-chan int, count int) int {
	return 0
}

// TODO: ForwardOne должен получить одно значение из in и отправить его в out.
// Если in закрыт, ничего не отправлять и вернуть false; иначе вернуть true.
func ForwardOne(in <-chan string, out chan<- string) bool {
	return false
}

// TODO: ReceiveOne должен получить одно значение с comma-ok.
// Для закрытого канала вернуть 0, false; ноль из открытого канала — обычное значение.
func ReceiveOne(values <-chan int) (int, bool) {
	return 0, false
}

// TODO: RelayN должен переслать не более count значений из in в out и вернуть
// количество пересланных значений. Закрытый in завершает работу раньше.
func RelayN(in <-chan int, out chan<- int, count int) int {
	return 0
}

func Example() string {
	pair := SendPair(10, 20)
	first, second := <-pair, <-pair
	sumInput := make(chan int, 3)
	sumInput <- 1
	sumInput <- 2
	sumInput <- 3
	out := make(chan string, 1)
	in := make(chan string, 1)
	in <- "go"
	forwarded := ForwardOne(in, out)
	text := ""
	if forwarded {
		text = <-out
	}
	return fmt.Sprintf("pair=%d,%d\nsum=%d\nforward=%t:%s", first, second, SumN(sumInput, 3), forwarded, text)
}

package selectflow

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrTimeout = errors.New("timeout")
	ErrClosed  = errors.New("channel closed")
)

// TODO: FirstReady должен выбрать готовый канал и вернуть value, source, true.
// source равен "left" или "right". Закрытый выбранный канал даёт "", source, false.
func FirstReady(left, right <-chan string) (string, string, bool) {
	return "", "", false
}

// TODO: TrySend должен неблокирующе отправить value.
// Если канал не готов принять значение, вернуть false через default.
func TrySend(out chan<- int, value int) bool {
	return false
}

// TODO: ReceiveOrTimeout должен вернуть значение из in, ErrClosed для закрытого in
// или ErrTimeout, когда первым готов timeout.
func ReceiveOrTimeout(in <-chan string, timeout <-chan time.Time) (string, error) {
	return "", nil
}

func Example() string {
	left := make(chan string, 1)
	left <- "left-value"
	value, source, ok := FirstReady(left, nil)
	out := make(chan int, 1)
	sent := TrySend(out, 7)
	ready := make(chan string, 1)
	ready <- "ready"
	received, _ := ReceiveOrTimeout(ready, nil)
	return fmt.Sprintf("first=%s:%s:%t\nsend=%t\nreceive=%s", source, value, ok, sent, received)
}

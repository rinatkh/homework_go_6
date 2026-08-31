package channels

import (
	"testing"
	"time"
)

func TestSendPair(t *testing.T) {
	tests := [][2]int{{0, 0}, {1, 2}, {-1, 1}, {10, 20}, {42, 7}, {-10, -20}, {5, 5}, {100, 0}, {0, 100}, {999, -999}}
	for _, tt := range tests {
		values := SendPair(tt[0], tt[1])
		first, second := <-values, <-values
		if first != tt[0] || second != tt[1] {
			t.Fatalf("SendPair(%d, %d) = %d, %d", tt[0], tt[1], first, second)
		}
		select {
		case _, ok := <-values:
			if !ok {
				t.Fatal("SendPair must not close the channel")
			}
			t.Fatal("SendPair sent more than two values")
		default:
		}
	}
}

func TestSumN(t *testing.T) {
	tests := []struct {
		values []int
		count  int
		want   int
	}{
		{nil, 0, 0},
		{[]int{1}, -1, 0},
		{[]int{1}, 1, 1},
		{[]int{1, 2}, 2, 3},
		{[]int{-1, 1}, 2, 0},
		{[]int{1, 2, 3}, 2, 3},
		{[]int{5, 5, 5}, 3, 15},
		{[]int{-5, -5}, 2, -10},
		{[]int{0, 0, 1}, 3, 1},
		{[]int{10, 20, 30, 40}, 4, 100},
	}
	for _, tt := range tests {
		values := make(chan int, len(tt.values))
		for _, value := range tt.values {
			values <- value
		}
		if got := SumN(values, tt.count); got != tt.want {
			t.Fatalf("SumN(%v, %d) = %d, want %d", tt.values, tt.count, got, tt.want)
		}
	}
}

func TestForwardOne(t *testing.T) {
	for _, value := range []string{"", "a", "go", "ready", "рус", "two words", "0", "!", "line\n"} {
		in := make(chan string, 1)
		out := make(chan string, 1)
		in <- value
		if ok := ForwardOne(in, out); !ok {
			t.Fatalf("ForwardOne(%q) = false", value)
		}
		select {
		case got := <-out:
			if got != value {
				t.Fatalf("ForwardOne(%q) sent %q", value, got)
			}
		case <-time.After(time.Second):
			t.Fatal("ForwardOne did not send a value")
		}
	}

	closed := make(chan string)
	close(closed)
	out := make(chan string, 1)
	if ok := ForwardOne(closed, out); ok {
		t.Fatal("ForwardOne(closed) = true")
	}
	if len(out) != 0 {
		t.Fatal("ForwardOne must not send after closed input")
	}
}

func TestExample(t *testing.T) {
	if got, want := Example(), "pair=10,20\nsum=6\nforward=true:go"; got != want {
		t.Fatalf("Example() = %q, want %q", got, want)
	}
}

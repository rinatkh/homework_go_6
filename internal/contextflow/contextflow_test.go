package contextflow

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestCheck(t *testing.T) {
	contexts := make([]context.Context, 0, 10)
	for i := 0; i < 5; i++ {
		contexts = append(contexts, context.Background())
	}
	for i := 0; i < 5; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		contexts = append(contexts, ctx)
	}
	for index, ctx := range contexts {
		err := Check(ctx)
		if index < 5 && err != nil {
			t.Fatalf("Check(active) = %v", err)
		}
		if index >= 5 && !errors.Is(err, context.Canceled) {
			t.Fatalf("Check(canceled) = %v", err)
		}
	}
}

func TestWaitValue(t *testing.T) {
	for _, want := range []string{"", "a", "go", "ready", "рус", "two words"} {
		values := make(chan string, 1)
		values <- want
		got, err := WaitValue(context.Background(), values)
		if err != nil || got != want {
			t.Fatalf("WaitValue() = %q, %v; want %q, nil", got, err, want)
		}
	}
	for i := 0; i < 2; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if _, err := WaitValue(ctx, nil); !errors.Is(err, context.Canceled) {
			t.Fatalf("WaitValue(canceled) = %v", err)
		}
	}
	for i := 0; i < 2; i++ {
		closed := make(chan string)
		close(closed)
		if _, err := WaitValue(context.Background(), closed); !errors.Is(err, ErrClosed) {
			t.Fatalf("WaitValue(closed) = %v", err)
		}
	}
}

func TestProcess(t *testing.T) {
	tests := []struct {
		in   []int
		fn   func(int) int
		want []int
	}{
		{nil, func(v int) int { return v }, []int{}},
		{[]int{}, func(v int) int { return v }, []int{}},
		{[]int{0}, func(v int) int { return v * 2 }, []int{0}},
		{[]int{1}, func(v int) int { return v * 2 }, []int{2}},
		{[]int{-1}, func(v int) int { return v * 2 }, []int{-2}},
		{[]int{1, 2}, func(v int) int { return v + 10 }, []int{11, 12}},
		{[]int{-2, 0, 2}, func(v int) int { return v * v }, []int{4, 0, 4}},
		{[]int{1, 2, 3}, nil, []int{1, 2, 3}},
		{[]int{5, 5}, func(int) int { return 1 }, []int{1, 1}},
		{[]int{10, -10}, func(v int) int { return v / 2 }, []int{5, -5}},
	}
	for _, tt := range tests {
		got, err := Process(context.Background(), tt.in, tt.fn)
		if err != nil || !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("Process(%v) = %v, %v; want %v, nil", tt.in, got, err, tt.want)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	got, err := Process(ctx, []int{1, 2, 3, 4}, func(value int) int {
		if value == 2 {
			cancel()
		}
		return value * 2
	})
	if !errors.Is(err, context.Canceled) || !reflect.DeepEqual(got, []int{2, 4}) {
		t.Fatalf("Process(partial cancel) = %v, %v", got, err)
	}
}

func TestSend(t *testing.T) {
	values := []int{0, 1, -1, 10, 42, -100}
	for _, value := range values {
		out := make(chan int, 1)
		if err := Send(context.Background(), out, value); err != nil {
			t.Fatalf("Send(%d) error = %v", value, err)
		}
		if got := <-out; got != value {
			t.Fatalf("Send(%d) sent %d", value, got)
		}
	}

	for index := 0; index < 4; index++ {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if err := Send(ctx, nil, index); !errors.Is(err, context.Canceled) {
			t.Fatalf("Send(canceled) error = %v", err)
		}
	}
}

func TestCollectN(t *testing.T) {
	channel := func(values []int, closed bool) <-chan int {
		result := make(chan int, len(values))
		for _, value := range values {
			result <- value
		}
		if closed {
			close(result)
		}
		return result
	}

	tests := []struct {
		name    string
		ctx     context.Context
		values  <-chan int
		count   int
		want    []int
		wantErr error
	}{
		{"zero", context.Background(), nil, 0, []int{}, nil},
		{"negative", context.Background(), nil, -1, []int{}, nil},
		{"one", context.Background(), channel([]int{1}, false), 1, []int{1}, nil},
		{"two of three", context.Background(), channel([]int{1, 2, 3}, false), 2, []int{1, 2}, nil},
		{"all three", context.Background(), channel([]int{-1, 0, 1}, false), 3, []int{-1, 0, 1}, nil},
		{"closed early", context.Background(), channel([]int{5}, true), 2, []int{5}, ErrClosed},
		{"closed empty", context.Background(), channel(nil, true), 1, []int{}, ErrClosed},
		{"duplicates", context.Background(), channel([]int{7, 7}, false), 2, []int{7, 7}, nil},
	}

	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	tests = append(tests, struct {
		name    string
		ctx     context.Context
		values  <-chan int
		count   int
		want    []int
		wantErr error
	}{"canceled", canceled, nil, 1, []int{}, context.Canceled})

	deadline, cancelDeadline := context.WithDeadline(context.Background(), time.Unix(0, 0))
	defer cancelDeadline()
	tests = append(tests, struct {
		name    string
		ctx     context.Context
		values  <-chan int
		count   int
		want    []int
		wantErr error
	}{"deadline", deadline, nil, 1, []int{}, context.DeadlineExceeded})

	for _, tt := range tests {
		got, err := CollectN(tt.ctx, tt.values, tt.count)
		if !reflect.DeepEqual(got, tt.want) || !errors.Is(err, tt.wantErr) {
			t.Fatalf("CollectN(%s) = %v, %v; want %v, %v", tt.name, got, err, tt.want, tt.wantErr)
		}
	}
}

func TestExample(t *testing.T) {
	if got, want := Example(), "check=context canceled\nvalue=done\nprocess=[2 4 6]"; got != want {
		t.Fatalf("Example() = %q, want %q", got, want)
	}
}

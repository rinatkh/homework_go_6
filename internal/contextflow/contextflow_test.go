package contextflow

import (
	"context"
	"errors"
	"reflect"
	"testing"
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

func TestExample(t *testing.T) {
	if got, want := Example(), "check=context canceled\nvalue=done\nprocess=[2 4 6]"; got != want {
		t.Fatalf("Example() = %q, want %q", got, want)
	}
}

package racemutex

import "testing"

func TestCounterAdd(t *testing.T) {
	for _, tt := range []struct {
		start int
		delta int
		want  int
	}{
		{0, 0, 0},
		{0, 1, 1},
		{1, 1, 2},
		{1, -1, 0},
		{-5, 2, -3},
		{10, -3, 7},
		{100, 50, 150},
		{-100, -50, -150},
		{7, 3, 10},
		{999, 1, 1000},
	} {
		counter := &Counter{value: tt.start}
		counter.Add(tt.delta)
		if got := counter.value; got != tt.want {
			t.Fatalf("Counter.Add(%d): value = %d, want %d", tt.delta, got, tt.want)
		}
	}
	var counter *Counter
	counter.Add(1)
}

func TestCounterValue(t *testing.T) {
	for _, want := range []int{0, 1, -1, 10, 42, -100, 7, 999, 5, 12345} {
		counter := &Counter{value: want}
		if got := counter.Value(); got != want {
			t.Fatalf("Counter.Value() = %d, want %d", got, want)
		}
	}
	var counter *Counter
	if got := counter.Value(); got != 0 {
		t.Fatalf("nil Counter.Value() = %d, want 0", got)
	}
}

func TestParallelAdd(t *testing.T) {
	tests := []struct {
		workers    int
		increments int
		want       int
	}{
		{0, 10, 0},
		{1, 0, 0},
		{-1, 10, 0},
		{1, 1, 1},
		{1, 10, 10},
		{2, 5, 10},
		{5, 2, 10},
		{10, 10, 100},
		{4, 25, 100},
		{8, 125, 1000},
	}
	for _, tt := range tests {
		if got := ParallelAdd(tt.workers, tt.increments); got != tt.want {
			t.Fatalf("ParallelAdd(%d, %d) = %d, want %d", tt.workers, tt.increments, got, tt.want)
		}
	}
}

func TestCounterSwap(t *testing.T) {
	tests := []struct {
		start   int
		value   int
		wantOld int
	}{
		{0, 0, 0},
		{0, 1, 0},
		{1, 0, 1},
		{-1, 1, -1},
		{10, -3, 10},
		{-5, -10, -5},
		{42, 42, 42},
		{100, 999, 100},
		{-100, 100, -100},
		{7, 3, 7},
	}
	for _, tt := range tests {
		counter := &Counter{value: tt.start}
		if got := counter.Swap(tt.value); got != tt.wantOld || counter.value != tt.value {
			t.Fatalf("Counter.Swap(%d) = %d, value = %d; want %d, %d", tt.value, got, counter.value, tt.wantOld, tt.value)
		}
	}

	var counter *Counter
	if got := counter.Swap(10); got != 0 {
		t.Fatalf("nil Counter.Swap() = %d, want 0", got)
	}
}

func TestParallelDeltas(t *testing.T) {
	tests := []struct {
		deltas []int
		want   int
	}{
		{nil, 0},
		{[]int{}, 0},
		{[]int{0}, 0},
		{[]int{1}, 1},
		{[]int{-1}, -1},
		{[]int{1, 2, 3}, 6},
		{[]int{-2, 0, 2}, 0},
		{[]int{5, 5, 5, 5}, 20},
		{[]int{100, -50, -25}, 25},
		{[]int{1, -1, 2, -2, 3, -3, 10}, 10},
	}
	for _, tt := range tests {
		if got := ParallelDeltas(tt.deltas); got != tt.want {
			t.Fatalf("ParallelDeltas(%v) = %d, want %d", tt.deltas, got, tt.want)
		}
	}
}

func TestExample(t *testing.T) {
	if got, want := Example(), "counter=400"; got != want {
		t.Fatalf("Example() = %q, want %q", got, want)
	}
}

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

func TestExample(t *testing.T) {
	if got, want := Example(), "counter=400"; got != want {
		t.Fatalf("Example() = %q, want %q", got, want)
	}
}

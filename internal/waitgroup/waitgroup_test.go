package waitgroup

import (
	"reflect"
	"sync/atomic"
	"testing"
)

func TestRunAll(t *testing.T) {
	for taskCount := 0; taskCount < 10; taskCount++ {
		var calls atomic.Int32
		tasks := make([]func(), taskCount)
		for index := range tasks {
			if index == 5 {
				continue
			}
			tasks[index] = func() { calls.Add(1) }
		}
		RunAll(tasks)
		want := int32(taskCount)
		if taskCount > 5 {
			want--
		}
		if got := calls.Load(); got != want {
			t.Fatalf("RunAll(%d tasks) calls = %d, want %d", taskCount, got, want)
		}
	}
}

func TestSquares(t *testing.T) {
	tests := []struct {
		in   []int
		want []int
	}{
		{nil, []int{}},
		{[]int{}, []int{}},
		{[]int{0}, []int{0}},
		{[]int{1}, []int{1}},
		{[]int{-2}, []int{4}},
		{[]int{2, 3}, []int{4, 9}},
		{[]int{-3, 0, 3}, []int{9, 0, 9}},
		{[]int{1, 2, 3, 4}, []int{1, 4, 9, 16}},
		{[]int{10, -10}, []int{100, 100}},
		{[]int{5, 5, 5}, []int{25, 25, 25}},
	}
	for _, tt := range tests {
		if got := Squares(tt.in); !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("Squares(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestApplyAll(t *testing.T) {
	tests := []struct {
		in   []int
		fn   func(int) int
		want []int
	}{
		{nil, func(v int) int { return v }, []int{}},
		{[]int{}, func(v int) int { return v }, []int{}},
		{[]int{1}, func(v int) int { return v * 2 }, []int{2}},
		{[]int{-1}, func(v int) int { return v * 2 }, []int{-2}},
		{[]int{1, 2}, func(v int) int { return v + 10 }, []int{11, 12}},
		{[]int{3, 2, 1}, func(v int) int { return v * v }, []int{9, 4, 1}},
		{[]int{0, 5}, func(v int) int { return -v }, []int{0, -5}},
		{[]int{7, 7}, func(int) int { return 1 }, []int{1, 1}},
		{[]int{1, 2, 3}, nil, []int{1, 2, 3}},
		{[]int{-10, 10}, func(v int) int { return v / 2 }, []int{-5, 5}},
	}
	for _, tt := range tests {
		got := ApplyAll(tt.in, tt.fn)
		if !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("ApplyAll(%v) = %v, want %v", tt.in, got, tt.want)
		}
		if len(got) > 0 && len(tt.in) > 0 {
			before := tt.in[0]
			got[0]++
			if tt.in[0] != before {
				t.Fatal("ApplyAll must return an independent slice")
			}
		}
	}
}

func TestExample(t *testing.T) {
	if got, want := Example(), "tasks=2\nsquares=[4 9]\napplied=[11 12]"; got != want {
		t.Fatalf("Example() = %q, want %q", got, want)
	}
}

package closerange

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

func collectWithTimeout(t *testing.T, values <-chan int) []int {
	t.Helper()
	done := make(chan []int, 1)
	go func() {
		done <- Collect(values)
	}()
	select {
	case result := <-done:
		return result
	case <-time.After(time.Second):
		t.Fatal("channel was not closed")
		return nil
	}
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		limit int
		want  []int
	}{
		{-2, []int{}},
		{-1, []int{}},
		{0, []int{}},
		{1, []int{0}},
		{2, []int{0, 1}},
		{3, []int{0, 1, 2}},
		{4, []int{0, 1, 2, 3}},
		{5, []int{0, 1, 2, 3, 4}},
		{6, []int{0, 1, 2, 3, 4, 5}},
		{10, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
	}
	for _, tt := range tests {
		if got := collectWithTimeout(t, Generate(tt.limit)); !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("Generate(%d) = %v, want %v", tt.limit, got, tt.want)
		}
	}
}

func TestCollect(t *testing.T) {
	tests := [][]int{
		nil,
		{},
		{0},
		{1},
		{-1},
		{1, 2},
		{-1, 0, 1},
		{1, 1, 1},
		{10, 20, 30},
		{1, 2, 3, 4, 5},
	}
	for _, values := range tests {
		if values == nil {
			if got := Collect(nil); !reflect.DeepEqual(got, []int{}) {
				t.Fatalf("Collect(nil) = %#v", got)
			}
			continue
		}
		input := make(chan int, len(values))
		for _, value := range values {
			input <- value
		}
		close(input)
		if got := Collect(input); !reflect.DeepEqual(got, values) {
			t.Fatalf("Collect(%v) = %v", values, got)
		}
	}
}

func TestFilterEven(t *testing.T) {
	tests := []struct {
		in   []int
		want []int
	}{
		{nil, []int{}},
		{[]int{}, []int{}},
		{[]int{1}, []int{}},
		{[]int{2}, []int{2}},
		{[]int{1, 2, 3, 4}, []int{2, 4}},
		{[]int{-2, -1, 0, 1, 2}, []int{-2, 0, 2}},
		{[]int{2, 2}, []int{2, 2}},
		{[]int{1, 3, 5}, []int{}},
		{[]int{10, 11, 12}, []int{10, 12}},
		{[]int{0, 1, 2, 3, 4, 5, 6}, []int{0, 2, 4, 6}},
	}
	for _, tt := range tests {
		if got := collectWithTimeout(t, FilterEven(tt.in)); !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("FilterEven(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestDouble(t *testing.T) {
	tests := []struct {
		in   []int
		want []int
	}{
		{nil, []int{}},
		{[]int{}, []int{}},
		{[]int{0}, []int{0}},
		{[]int{1}, []int{2}},
		{[]int{-1}, []int{-2}},
		{[]int{1, 2, 3}, []int{2, 4, 6}},
		{[]int{-2, 0, 2}, []int{-4, 0, 4}},
		{[]int{5, 5}, []int{10, 10}},
		{[]int{10, -10}, []int{20, -20}},
		{[]int{1, 2, 3, 4, 5}, []int{2, 4, 6, 8, 10}},
	}
	for _, tt := range tests {
		if got := collectWithTimeout(t, Double(tt.in)); !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("Double(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestMerge(t *testing.T) {
	makeInput := func(values []int, useNil bool) <-chan int {
		if useNil {
			return nil
		}
		input := make(chan int, len(values))
		for _, value := range values {
			input <- value
		}
		close(input)
		return input
	}
	tests := []struct {
		leftNil  bool
		rightNil bool
		left     []int
		right    []int
		want     []int
	}{
		{true, true, nil, nil, []int{}},
		{false, false, []int{}, []int{}, []int{}},
		{false, true, []int{1}, nil, []int{1}},
		{true, false, nil, []int{2}, []int{2}},
		{false, false, []int{1}, []int{2}, []int{1, 2}},
		{false, false, []int{1, 3}, []int{2, 4}, []int{1, 2, 3, 4}},
		{false, false, []int{-1, 0}, []int{1}, []int{-1, 0, 1}},
		{false, false, []int{5, 5}, []int{5}, []int{5, 5, 5}},
		{false, false, []int{10, 20, 30}, []int{}, []int{10, 20, 30}},
		{false, false, []int{100, -100}, []int{7, 8, 9}, []int{-100, 7, 8, 9, 100}},
	}
	for _, tt := range tests {
		got := collectWithTimeout(t, Merge(makeInput(tt.left, tt.leftNil), makeInput(tt.right, tt.rightNil)))
		sort.Ints(got)
		want := append([]int{}, tt.want...)
		sort.Ints(want)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Merge(%v, %v) = %v, want values %v", tt.left, tt.right, got, want)
		}
	}
}

func TestExample(t *testing.T) {
	if got, want := Example(), "generated=[0 1 2 3 4]\neven=[2 4 6]"; got != want {
		t.Fatalf("Example() = %q, want %q", got, want)
	}
}

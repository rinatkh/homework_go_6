package closerange

import (
	"reflect"
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

func TestExample(t *testing.T) {
	if got, want := Example(), "generated=[0 1 2 3 4]\neven=[2 4 6]"; got != want {
		t.Fatalf("Example() = %q, want %q", got, want)
	}
}

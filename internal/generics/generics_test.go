package generics

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name string
		got  any
		want any
	}{
		{"nil", Map([]int(nil), func(v int) string { return fmt.Sprint(v) }), []string{}},
		{"empty", Map([]string{}, func(v string) int { return len(v) }), []int{}},
		{"squares", Map([]int{1, 2, 3}, func(v int) int { return v * v }), []int{1, 4, 9}},
		{"negative", Map([]int{-2, 0, 2}, func(v int) int { return v + 1 }), []int{-1, 1, 3}},
		{"int to string", Map([]int{1, 2}, func(v int) string { return fmt.Sprintf("#%d", v) }), []string{"#1", "#2"}},
		{"string lengths", Map([]string{"go", "generics"}, func(v string) int { return len(v) }), []int{2, 8}},
		{"empty string", Map([]string{"", "x"}, func(v string) bool { return v == "" }), []bool{true, false}},
		{"booleans", Map([]bool{true, false}, func(v bool) string { return fmt.Sprint(v) }), []string{"true", "false"}},
		{"order", Map([]int{3, 1, 2}, func(v int) int { return v * 10 }), []int{30, 10, 20}},
		{"duplicates", Map([]string{"a", "a"}, func(v string) string { return v + v }), []string{"aa", "aa"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Fatalf("Map() = %#v, want %#v", tt.got, tt.want)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	type point struct{ X, Y int }
	tests := []struct {
		name string
		got  any
		want any
	}{
		{"nil", Unique([]int(nil)), []int{}},
		{"empty", Unique([]string{}), []string{}},
		{"one", Unique([]int{1}), []int{1}},
		{"no duplicates", Unique([]int{1, 2, 3}), []int{1, 2, 3}},
		{"duplicates", Unique([]int{2, 1, 2, 3, 1}), []int{2, 1, 3}},
		{"all same", Unique([]int{5, 5, 5}), []int{5}},
		{"strings", Unique([]string{"go", "go", "generic"}), []string{"go", "generic"}},
		{"empty string", Unique([]string{"", "x", ""}), []string{"", "x"}},
		{"bool", Unique([]bool{true, false, true}), []bool{true, false}},
		{"struct", Unique([]point{{1, 2}, {1, 2}, {3, 4}}), []point{{1, 2}, {3, 4}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Fatalf("Unique() = %#v, want %#v", tt.got, tt.want)
			}
		})
	}
}

func TestDrain(t *testing.T) {
	drainInts := func(values ...int) []int {
		ch := make(chan int, len(values))
		for _, value := range values {
			ch <- value
		}
		close(ch)
		return Drain(ch)
	}
	drainStrings := func(values ...string) []string {
		ch := make(chan string, len(values))
		for _, value := range values {
			ch <- value
		}
		close(ch)
		return Drain(ch)
	}

	tests := []struct {
		name string
		got  any
		want any
	}{
		{"nil", Drain((<-chan int)(nil)), []int{}},
		{"empty ints", drainInts(), []int{}},
		{"one int", drainInts(1), []int{1}},
		{"many ints", drainInts(1, 2, 3), []int{1, 2, 3}},
		{"negative ints", drainInts(-2, 0, 2), []int{-2, 0, 2}},
		{"empty strings", drainStrings(), []string{}},
		{"one string", drainStrings("go"), []string{"go"}},
		{"many strings", drainStrings("go", "generic"), []string{"go", "generic"}},
		{"empty value", drainStrings(""), []string{""}},
		{"duplicates", drainStrings("x", "x"), []string{"x", "x"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Fatalf("Drain() = %#v, want %#v", tt.got, tt.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name string
		got  any
		want any
	}{
		{"nil", Filter([]int(nil), func(v int) bool { return v > 0 }), []int{}},
		{"empty", Filter([]string{}, func(v string) bool { return v != "" }), []string{}},
		{"all", Filter([]int{1, 2, 3}, func(v int) bool { return v > 0 }), []int{1, 2, 3}},
		{"none", Filter([]int{1, 2, 3}, func(v int) bool { return v < 0 }), []int{}},
		{"even", Filter([]int{-2, -1, 0, 1, 2}, func(v int) bool { return v%2 == 0 }), []int{-2, 0, 2}},
		{"strings", Filter([]string{"", "go", "generic"}, func(v string) bool { return len(v) >= 2 }), []string{"go", "generic"}},
		{"bool", Filter([]bool{true, false, true}, func(v bool) bool { return v }), []bool{true, true}},
		{"duplicates", Filter([]string{"x", "x", "y"}, func(v string) bool { return v == "x" }), []string{"x", "x"}},
		{"keep order", Filter([]int{3, 1, 2}, func(v int) bool { return v != 1 }), []int{3, 2}},
		{"nil predicate", Filter([]int{1, 2, 3}, nil), []int{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Fatalf("Filter() = %#v, want %#v", tt.got, tt.want)
			}
		})
	}

	input := []int{1, 2}
	got := Filter(input, nil)
	if len(got) > 0 {
		got[0] = 99
	}
	if input[0] != 1 {
		t.Fatal("Filter with nil predicate must return an independent slice")
	}
}

func TestIndexOf(t *testing.T) {
	type point struct{ X, Y int }
	tests := []struct {
		name string
		got  int
		want int
	}{
		{"nil", IndexOf([]int(nil), 1), -1},
		{"empty", IndexOf([]string{}, "go"), -1},
		{"first", IndexOf([]int{1, 2, 3}, 1), 0},
		{"middle", IndexOf([]int{1, 2, 3}, 2), 1},
		{"last", IndexOf([]int{1, 2, 3}, 3), 2},
		{"missing", IndexOf([]int{1, 2, 3}, 4), -1},
		{"first duplicate", IndexOf([]string{"go", "x", "go"}, "go"), 0},
		{"empty string", IndexOf([]string{"go", ""}, ""), 1},
		{"bool", IndexOf([]bool{false, true}, true), 1},
		{"struct", IndexOf([]point{{1, 2}, {3, 4}}, point{3, 4}), 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("IndexOf() = %d, want %d", tt.got, tt.want)
			}
		})
	}
}

func TestExample(t *testing.T) {
	if got, want := Example(), "map=[item-1 item-2 item-3]\nunique=[2 1 3]\ndrain=[go generic]"; got != want {
		t.Fatalf("Example() = %q, want %q", got, want)
	}
}

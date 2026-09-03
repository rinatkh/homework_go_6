package goroutines

import (
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func receiveInt(t *testing.T, values <-chan int) int {
	t.Helper()
	select {
	case value := <-values:
		return value
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for int result")
		return 0
	}
}

func receiveString(t *testing.T, values <-chan string) string {
	t.Helper()
	select {
	case value := <-values:
		return value
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for string result")
		return ""
	}
}

func waitDone(t *testing.T, done <-chan struct{}) {
	t.Helper()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for job")
	}
}

func TestComputeAsync(t *testing.T) {
	tests := []struct {
		value int
		fn    func(int) int
		want  int
	}{
		{0, func(v int) int { return v * v }, 0},
		{1, func(v int) int { return v * v }, 1},
		{2, func(v int) int { return v * v }, 4},
		{-3, func(v int) int { return v * v }, 9},
		{10, func(v int) int { return v + 1 }, 11},
		{-10, func(v int) int { return v - 1 }, -11},
		{7, func(v int) int { return v * 2 }, 14},
		{9, func(v int) int { return -v }, -9},
		{42, func(v int) int { return v }, 42},
		{5, nil, 5},
	}
	for _, tt := range tests {
		if got := receiveInt(t, ComputeAsync(tt.value, tt.fn)); got != tt.want {
			t.Fatalf("ComputeAsync(%d) = %d, want %d", tt.value, got, tt.want)
		}
	}
}

func TestStringAsync(t *testing.T) {
	tests := []struct {
		value string
		fn    func(string) string
		want  string
	}{
		{"", strings.ToUpper, ""},
		{"go", strings.ToUpper, "GO"},
		{"Go", strings.ToLower, "go"},
		{" go ", strings.TrimSpace, "go"},
		{"рус", strings.ToUpper, "РУС"},
		{"a", func(v string) string { return v + v }, "aa"},
		{"x", func(string) string { return "fixed" }, "fixed"},
		{"line\n", strings.TrimSpace, "line"},
		{"123", strings.ToUpper, "123"},
		{"keep", nil, "keep"},
	}
	for _, tt := range tests {
		if got := receiveString(t, StringAsync(tt.value, tt.fn)); got != tt.want {
			t.Fatalf("StringAsync(%q) = %q, want %q", tt.value, got, tt.want)
		}
	}
}

func TestRunAsync(t *testing.T) {
	tests := []bool{true, true, false, true, false, true, true, false, true, false}
	for _, hasJob := range tests {
		var calls atomic.Int32
		var job func()
		if hasJob {
			job = func() { calls.Add(1) }
		}
		waitDone(t, RunAsync(job))
		var want int32
		if hasJob {
			want = 1
		}
		if got := calls.Load(); got != want {
			t.Fatalf("RunAsync(hasJob=%t) calls = %d, want %d", hasJob, got, want)
		}
	}
}

func TestSumAsync(t *testing.T) {
	tests := []struct {
		values []int
		want   int
	}{
		{nil, 0},
		{[]int{}, 0},
		{[]int{0}, 0},
		{[]int{1}, 1},
		{[]int{-1}, -1},
		{[]int{1, 2, 3}, 6},
		{[]int{-2, 0, 2}, 0},
		{[]int{5, 5, 5}, 15},
		{[]int{10, -3, -2}, 5},
		{[]int{100, 200, 300, 400}, 1000},
	}
	for _, tt := range tests {
		if got := receiveInt(t, SumAsync(tt.values)); got != tt.want {
			t.Fatalf("SumAsync(%v) = %d, want %d", tt.values, got, tt.want)
		}
	}
}

func TestCountAsync(t *testing.T) {
	tests := []struct {
		values    []string
		predicate func(string) bool
		want      int
	}{
		{nil, nil, 0},
		{[]string{}, func(string) bool { return true }, 0},
		{[]string{"go"}, nil, 1},
		{[]string{"", "go"}, func(v string) bool { return v == "" }, 1},
		{[]string{"a", "bb", "ccc"}, func(v string) bool { return len(v) >= 2 }, 2},
		{[]string{"go", "rust", "go"}, func(v string) bool { return v == "go" }, 2},
		{[]string{"A", "b", "C"}, func(v string) bool { return strings.ToUpper(v) == v }, 2},
		{[]string{"1", "22", "333"}, func(v string) bool { return len(v)%2 == 1 }, 2},
		{[]string{"x", "x", "x"}, func(string) bool { return false }, 0},
		{[]string{"рус", "go", "язык"}, func(v string) bool { return strings.Contains(v, "я") }, 1},
	}
	for _, tt := range tests {
		if got := receiveInt(t, CountAsync(tt.values, tt.predicate)); got != tt.want {
			t.Fatalf("CountAsync(%v) = %d, want %d", tt.values, got, tt.want)
		}
	}
}

func TestCountAsyncStartsWorker(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	result := CountAsync([]string{"go"}, func(string) bool {
		close(started)
		<-release
		return true
	})

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("CountAsync не запустил функцию проверки в отдельной горутине")
	}

	select {
	case <-result:
		t.Fatal("CountAsync returned a result before predicate completed")
	default:
	}

	close(release)
	if got := receiveInt(t, result); got != 1 {
		t.Fatalf("CountAsync() = %d, want 1", got)
	}
}

func TestExample(t *testing.T) {
	if got, want := Example(), "compute=16\nstring=GO\njob=true"; got != want {
		t.Fatalf("Example() = %q, want %q", got, want)
	}
}

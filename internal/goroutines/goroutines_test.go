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

func TestExample(t *testing.T) {
	if got, want := Example(), "compute=16\nstring=GO\njob=true"; got != want {
		t.Fatalf("Example() = %q, want %q", got, want)
	}
}

package selectflow

import (
	"errors"
	"testing"
	"time"
)

func readyString(value string) <-chan string {
	values := make(chan string, 1)
	values <- value
	return values
}

func closedString() <-chan string {
	values := make(chan string)
	close(values)
	return values
}

func TestFirstReady(t *testing.T) {
	tests := []struct {
		left       <-chan string
		right      <-chan string
		wantValue  string
		wantSource string
		wantOK     bool
		either     bool
	}{
		{readyString("a"), nil, "a", "left", true, false},
		{nil, readyString("b"), "b", "right", true, false},
		{readyString(""), nil, "", "left", true, false},
		{nil, readyString("рус"), "рус", "right", true, false},
		{readyString("left"), readyString("right"), "", "", true, true},
		{closedString(), nil, "", "left", false, false},
		{nil, closedString(), "", "right", false, false},
		{readyString("two words"), nil, "two words", "left", true, false},
		{nil, readyString("line\n"), "line\n", "right", true, false},
		{readyString("last"), nil, "last", "left", true, false},
	}
	for _, tt := range tests {
		value, source, ok := FirstReady(tt.left, tt.right)
		if tt.either {
			valid := ok && ((source == "left" && value == "left") || (source == "right" && value == "right"))
			if !valid {
				t.Fatalf("FirstReady(both) = %q, %q, %t", value, source, ok)
			}
			continue
		}
		if value != tt.wantValue || source != tt.wantSource || ok != tt.wantOK {
			t.Fatalf("FirstReady() = %q, %q, %t; want %q, %q, %t", value, source, ok, tt.wantValue, tt.wantSource, tt.wantOK)
		}
	}
}

func TestTrySend(t *testing.T) {
	for capacity := 0; capacity < 10; capacity++ {
		out := make(chan int, capacity)
		got := TrySend(out, capacity)
		want := capacity > 0
		if got != want {
			t.Fatalf("TrySend(capacity=%d) = %t, want %t", capacity, got, want)
		}
		if got && <-out != capacity {
			t.Fatalf("TrySend(capacity=%d) sent wrong value", capacity)
		}
	}
}

func TestReceiveOrTimeout(t *testing.T) {
	for _, want := range []string{"", "a", "go", "ready", "рус", "two words"} {
		got, err := ReceiveOrTimeout(readyString(want), nil)
		if err != nil || got != want {
			t.Fatalf("ReceiveOrTimeout() = %q, %v; want %q, nil", got, err, want)
		}
	}

	for i := 0; i < 2; i++ {
		timeout := make(chan time.Time)
		close(timeout)
		if _, err := ReceiveOrTimeout(nil, timeout); !errors.Is(err, ErrTimeout) {
			t.Fatalf("ReceiveOrTimeout timeout = %v", err)
		}
	}

	for i := 0; i < 2; i++ {
		closed := make(chan string)
		close(closed)
		if _, err := ReceiveOrTimeout(closed, nil); !errors.Is(err, ErrClosed) {
			t.Fatalf("ReceiveOrTimeout closed = %v", err)
		}
	}
}

func TestExample(t *testing.T) {
	if got, want := Example(), "first=left:left-value:true\nsend=true\nreceive=ready"; got != want {
		t.Fatalf("Example() = %q, want %q", got, want)
	}
}

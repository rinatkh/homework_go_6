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

func TestTryReceive(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		buffered bool
		closed   bool
		useNil   bool
		want     int
		wantOK   bool
	}{
		{"nil", 0, false, false, true, 0, false},
		{"not ready", 0, false, false, false, 0, false},
		{"closed", 0, false, true, false, 0, false},
		{"zero", 0, true, false, false, 0, true},
		{"one", 1, true, false, false, 1, true},
		{"negative", -1, true, false, false, -1, true},
		{"forty two", 42, true, false, false, 42, true},
		{"large", 1000, true, false, false, 1000, true},
		{"minus large", -1000, true, false, false, -1000, true},
		{"seven", 7, true, false, false, 7, true},
	}
	for _, tt := range tests {
		var in <-chan int
		if !tt.useNil {
			values := make(chan int, 1)
			if tt.buffered {
				values <- tt.value
			}
			if tt.closed {
				close(values)
			}
			in = values
		}
		got, ok := TryReceive(in)
		if got != tt.want || ok != tt.wantOK {
			t.Fatalf("TryReceive(%s) = %d, %t; want %d, %t", tt.name, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestPrimaryOrFallback(t *testing.T) {
	tests := []struct {
		name       string
		primary    <-chan string
		fallback   <-chan string
		wantValue  string
		wantSource string
		wantOK     bool
	}{
		{"primary only", readyString("p"), nil, "p", "primary", true},
		{"fallback only", nil, readyString("f"), "f", "fallback", true},
		{"both ready prefer primary", readyString("p"), readyString("f"), "p", "primary", true},
		{"closed primary uses fallback", closedString(), readyString("f"), "f", "fallback", true},
		{"closed fallback uses primary", readyString("p"), closedString(), "p", "primary", true},
		{"both closed", closedString(), closedString(), "", "", false},
		{"both nil", nil, nil, "", "", false},
		{"closed primary only", closedString(), nil, "", "", false},
		{"closed fallback only", nil, closedString(), "", "", false},
		{"empty primary value", readyString(""), readyString("f"), "", "primary", true},
	}
	for _, tt := range tests {
		done := make(chan struct{})
		var gotValue, gotSource string
		var gotOK bool
		go func() {
			gotValue, gotSource, gotOK = PrimaryOrFallback(tt.primary, tt.fallback)
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatalf("PrimaryOrFallback(%s) did not finish", tt.name)
		}
		if gotValue != tt.wantValue || gotSource != tt.wantSource || gotOK != tt.wantOK {
			t.Fatalf("PrimaryOrFallback(%s) = %q, %q, %t; want %q, %q, %t", tt.name, gotValue, gotSource, gotOK, tt.wantValue, tt.wantSource, tt.wantOK)
		}
	}
}

func TestExample(t *testing.T) {
	if got, want := Example(), "first=left:left-value:true\nsend=true\nreceive=ready"; got != want {
		t.Fatalf("Example() = %q, want %q", got, want)
	}
}

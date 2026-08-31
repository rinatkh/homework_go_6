package integration

import (
	"os/exec"
	"strings"
	"testing"
)

func TestCommands(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"../../cmd/01_goroutines", "compute=16\nstring=GO\njob=true"},
		{"../../cmd/02_channels", "pair=10,20\nsum=6\nforward=true:go"},
		{"../../cmd/03_waitgroup", "tasks=2\nsquares=[4 9]\napplied=[11 12]"},
		{"../../cmd/04_race_mutex", "counter=400"},
		{"../../cmd/05_close_range", "generated=[0 1 2 3 4]\neven=[2 4 6]"},
		{"../../cmd/06_select", "first=left:left-value:true\nsend=true\nreceive=ready"},
		{"../../cmd/07_context", "check=context canceled\nvalue=done\nprocess=[2 4 6]"},
		{"../../cmd/08_generics", "map=[item-1 item-2 item-3]\nunique=[2 1 3]\ndrain=[go generic]"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			command := exec.Command("go", "run", tt.path)
			output, err := command.CombinedOutput()
			if err != nil {
				t.Fatalf("go run %s: %v\n%s", tt.path, err, output)
			}
			if got := strings.TrimSpace(string(output)); got != tt.want {
				t.Fatalf("go run %s = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

package commands

import (
	"testing"
)

func TestHandleINF(t *testing.T) {
	tests := []struct {
		arguments string
		expected  string
		ok        bool
	}{
		{"1", "INF 1 MD5\r\n", true},
		{"", "", false},
	}

	for _, tt := range tests {
		c := make(chan string, 1)

		err := HandleINF(c, tt.arguments)

		if (err == nil) != tt.ok {
			t.Errorf("HandleINF(%q) = %v, want %v", tt.arguments, err == nil, tt.ok)
			t.FailNow()
		}

		if tt.expected != "" {
			if got := <-c; got != tt.expected {
				t.Errorf("HandleINF(%q) = %q, want %q", tt.arguments, got, tt.expected)
			}
		}
	}
}

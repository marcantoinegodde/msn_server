package commands

import (
	"testing"
)

func TestHandleINF(t *testing.T) {
	tests := []struct {
		arguments string
		expected  string
	}{
		{"1", "INF 1 MD5\r\n"},
		{"", ""},
	}

	for _, tt := range tests {
		mock := &mockConn{}
		HandleINF(mock, tt.arguments)

		if got := mock.buffer.String(); got != tt.expected {
			t.Errorf("HandleINF(%q) = %q, want %q", tt.arguments, got, tt.expected)
		}
	}
}

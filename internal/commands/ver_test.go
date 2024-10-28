package commands

import (
	"testing"
)

func TestHandleVER(t *testing.T) {
	tests := []struct {
		arguments string
		expected  string
		ok        bool
	}{
		{"1 MSNP2 CVR0\r\n", "VER 1 MSNP2 CVR0\r\n", true},
		{"100 MSNP7 MSNP6 MSNP5 MSNP4 MSNP2 CVR0\r\n", "VER 100 MSNP2 CVR0\r\n", true},
		{"1000 MYPROTOCOL\r\n", "VER 1000 0\r\n", false},
		{"MSNP2 CVR0\r\n", "", false},
		{"1 CVR0\r\n", "VER 1 0\r\n", false},
	}

	for _, tt := range tests {
		mock := &mockConn{}
		err := HandleVER(mock, tt.arguments)

		if got := mock.buffer.String(); got != tt.expected {
			t.Errorf("HandleVER(%q) = %q, want %q", tt.arguments, got, tt.expected)
		}
		if (err == nil) != tt.ok {
			t.Errorf("HandleVER(%q) = %v, want %v", tt.arguments, err == nil, tt.ok)
		}
	}
}

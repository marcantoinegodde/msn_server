package commands

import (
	"testing"
)

func TestHandleINF(t *testing.T) {
	tests := []struct {
		transactionID string
		expected      string
	}{
		{"1", "INF 1 MD5\r\n"},
		{"100", "INF 100 MD5\r\n"},
	}

	for _, tt := range tests {
		mock := &mockConn{}
		HandleINF(mock, tt.transactionID)

		if got := mock.buffer.String(); got != tt.expected {
			t.Errorf("HandleINF(%q) = %q, want %q", tt.transactionID, got, tt.expected)
		}
	}
}

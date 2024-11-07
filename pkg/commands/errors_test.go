package commands

import (
	"testing"
)

func TestSendError(t *testing.T) {
	tests := []struct {
		transactionID string
		errorCode     int
		expected      string
	}{
		{"123", ERR_SYNTAX_ERROR, "200 123\r\n"},
		{"456", ERR_INVALID_PARAMETER, "201 456\r\n"},
		{"789", ERR_INVALID_USER, "205 789\r\n"},
	}

	for _, tt := range tests {
		conn := &mockConn{}
		SendError(conn, tt.transactionID, tt.errorCode)

		if conn.buffer.String() != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, conn.buffer.String())
		}
	}
}

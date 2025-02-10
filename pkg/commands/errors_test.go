package commands

import (
	"msnserver/pkg/clients"
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
		c := &clients.Client{
			SendChan: make(chan string),
		}

		go SendError(c, tt.transactionID, tt.errorCode)

		if got := <-c.SendChan; got != tt.expected {
			t.Errorf("expected %q, got %q", tt.expected, got)
		}
	}
}

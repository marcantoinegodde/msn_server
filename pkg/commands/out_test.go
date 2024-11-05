package commands

import (
	"testing"
)

func TestHandleOUT(t *testing.T) {
	expected := "OUT\r\n"

	mock := &mockConn{}
	HandleOUT(mock)

	if got := mock.buffer.String(); got != expected {
		t.Errorf("HandleOUT() = %q, want %q", got, expected)
	}
}

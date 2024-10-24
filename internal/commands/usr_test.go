package commands

import (
	"msnserver/config"
	"testing"
)

func TestHandleUSRDispatch(t *testing.T) {
	config.Config.DispatchServer.ServerAddr = "207.46.104.20"
	config.Config.DispatchServer.ServerPort = "1863"
	config.Config.DispatchServer.NotificationServerAddr = "207.46.106.145"
	config.Config.DispatchServer.NotificationServerPort = "1863"

	tests := []struct {
		transactionID string
		expected      string
	}{
		{"1", "XFR 1 NS 207.46.106.145:1863 0 207.46.104.20:1863\r\n"},
		{"100", "XFR 100 NS 207.46.106.145:1863 0 207.46.104.20:1863\r\n"},
	}

	for _, tt := range tests {
		mock := &mockConn{}
		HandleUSRDispatch(mock, tt.transactionID)

		if got := mock.buffer.String(); got != tt.expected {
			t.Errorf("HandleUSRDispatch(%q) = %q, want %q", tt.transactionID, got, tt.expected)
		}
	}
}

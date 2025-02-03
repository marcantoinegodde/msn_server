package commands

import (
	"msnserver/pkg/clients"
	"testing"
)

func TestHandleOUT(t *testing.T) {
	tests := []struct {
		name     string
		reason   string
		expected string
	}{
		{"Valid reason OTH", "OTH", "OUT OTH\r\n"},
		{"Valid reason SSD", "SSD", "OUT SSD\r\n"},
		{"Blanck reading", "", "OUT\r\n"},
		{"Invalid reason", "INVALID", "OUT\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &clients.Client{
				SendChan: make(chan string, 1),
			}

			go HandleOUT(c, tt.reason)

			if got := <-c.SendChan; got != tt.expected {
				t.Errorf("HandleOUT() = %q, want %q", got, tt.expected)
			}
		})
	}
}

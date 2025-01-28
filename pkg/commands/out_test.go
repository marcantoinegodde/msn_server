package commands

import (
	"msnserver/pkg/clients"
	"testing"
)

func TestHandleOUT(t *testing.T) {
	c := &clients.Client{
		SendChan: make(chan string),
	}
	expected := "OUT\r\n"

	go HandleOUT(c)

	if got := <-c.SendChan; got != expected {
		t.Errorf("HandleOUT() = %q, want %q", got, expected)
	}
}

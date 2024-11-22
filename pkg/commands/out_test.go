package commands

import (
	"testing"
)

func TestHandleOUT(t *testing.T) {
	c := make(chan string)
	expected := "OUT\r\n"

	go HandleOUT(c)

	if got := <-c; got != expected {
		t.Errorf("HandleOUT() = %q, want %q", got, expected)
	}
}

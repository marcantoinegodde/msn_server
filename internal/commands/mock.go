package commands

import (
	"bytes"
	"net"
)

type mockConn struct {
	net.Conn
	buffer bytes.Buffer
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	return m.buffer.Write(b)
}

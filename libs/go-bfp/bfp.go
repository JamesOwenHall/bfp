package bfp

import (
	"encoding/json"
	"errors"
	"github.com/JamesOwenHall/BruteForceProtection/core/message-server"
	"net"
)

const (
	UnixType = "unix"
	TcpType  = "tcp"
)

var ConnectionError = errors.New("a connection error occurred")

type Bfp struct {
	Type string
	Addr string
}

func (b *Bfp) Hit(direction string, value interface{}) (bool, error) {
	conn, err := net.Dial(b.Type, b.Addr)
	if err != nil {
		return false, ConnectionError
	}
	defer conn.Close()

	message := server.Request{
		Direction: direction,
		Value:     value,
	}

	enc := json.NewEncoder(conn)
	err = enc.Encode(message)
	if err != nil {
		return false, ConnectionError
	}

	buf := []byte{0}
	_, err = conn.Read(buf)
	if err != nil {
		return false, ConnectionError
	}

	if buf[0] == 't' {
		return true, nil
	} else {
		return false, nil
	}
}

// Package bfp provides the API for communicating with an instance of BFP.
package bfp

import (
	"encoding/json"
	"errors"
	"github.com/JamesOwenHall/bfp/core/message-server"
	"net"
)

// The permitted connection types.
const (
	UnixType = "unix"
	TcpType  = "tcp"
)

var ConnectionError = errors.New("a connection error occurred")

// Bfp is a struct used to communicate with an instance of BFP.
type Bfp struct {
	Type string
	Addr string
}

// Hit registers the use of a value.  It returns true if the request is deemed
// safe and false otherwise.  If the error is not nil, this function always
// returns false.
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

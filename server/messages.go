package server

import (
	"encoding/json"
	"io"
)

type Request struct {
	Direction string
	Value     interface{}
}

func ReadRequest(r io.Reader) (*Request, error) {
	decoder := json.NewDecoder(r)
	result := new(Request)
	err := decoder.Decode(result)
	return result, err
}

type Response struct {
	Valid bool
}

func (r *Response) Write(w io.Writer) error {
	encoder := json.NewEncoder(w)
	err := encoder.Encode(*r)
	return err
}

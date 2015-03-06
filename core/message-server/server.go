// Package server implements a server following its own custom protocol.  The
// protocol works as follows.
//
//     1. A client connects and sends an object in JSON format matching the
//     structure of the Request type.
//
//     2. The server returns a single character. "t" means that the request is
//     valid.  "f" means that the request is invalid (either an attack or an
//     error).
//
//     3. The client disconnects or goes again from step 1.
package server

import (
	"log"
	"net"
)

// Server is a server that interprets requests according to the protocol.
type Server struct {
	Routes   map[string]func(interface{}) bool
	listener net.Listener
}

// New returns an initialized *Server.
func New() *Server {
	return &Server{
		Routes: make(map[string]func(interface{}) bool),
	}
}

// Blocks and listens for requests.
func (s *Server) ListenAndServe(typ, addr string) error {
	// Start listening.
	var err error
	s.listener, err = net.Listen(typ, addr)
	if err != nil {
		return err
	}

	// Accept requests.
	go s.acceptRequests()
	return nil
}

func (s *Server) acceptRequests() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("connection error:", err)
			continue
		}

		// For performance, we launch every handler in its own goroutine.
		go func(conn net.Conn) {
			for {
				request, err := ReadRequest(conn)
				if err != nil {
					conn.Close()
					return
				}

				// Check if the route exists.
				if handler, ok := s.Routes[request.Direction]; ok {
					response := handler(request.Value)
					if response {
						conn.Write([]byte("t"))
					} else {
						conn.Write([]byte("f"))
					}
				}
			}
		}(conn)
	}
}

// Close stops the server.
func (s *Server) Close() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	s.listener.Close()
}

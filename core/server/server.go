package server

import (
	"log"
	"net"
)

type Server struct {
	Routes   map[string]func(interface{}) Response
	listener net.Listener
}

func New() *Server {
	return &Server{
		Routes: make(map[string]func(interface{}) Response),
	}
}

func (s *Server) ListenAndServe(addr string) {
	// Start listening.
	var err error
	s.listener, err = net.Listen("unix", addr)
	if err != nil {
		panic(err)
	}

	// Accept requests.
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("connection error:", err)
			continue
		}

		// For performance, we launch every handler in its own goroutine.
		go func(conn net.Conn) {
			request, err := ReadRequest(conn)
			if err != nil {
				log.Println("Failed to read request:", err)
				return
			}

			// Check if the route exists.
			if handler, ok := s.Routes[request.Direction]; ok {
				response := handler(request.Value)
				response.Write(conn)
			}

			// Remember to close the connection
			conn.Close()
		}(conn)
	}
}

func (s *Server) Close() {
	s.listener.Close()
}

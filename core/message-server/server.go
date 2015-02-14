package server

import (
	"log"
	"net"
)

type Server struct {
	Routes   map[string]func(interface{}) bool
	listener net.Listener
}

func New() *Server {
	return &Server{
		Routes: make(map[string]func(interface{}) bool),
	}
}

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

func (s *Server) Close() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	s.listener.Close()
}

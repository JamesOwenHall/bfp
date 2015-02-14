package dashboard

import (
	"fmt"
	"net/http"
)

type Server struct {
	http.Server
	mux *http.ServeMux
}

func New(addr string) *Server {
	result := new(Server)
	result.mux = http.NewServeMux()
	result.Server = http.Server{
		Addr:    addr,
		Handler: result.mux,
	}

	result.setupRoutes()

	return result
}

func (s *Server) ListenAndServe() {
	go s.Server.ListenAndServe()
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<h1>Dashboard</h1>")
	})
}

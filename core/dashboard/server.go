// Package dashboard provides a web portal to see the status of the system.
package dashboard

import (
	"github.com/JamesOwenHall/bfp/core/config"
	"github.com/JamesOwenHall/bfp/core/hitcounter"
	"html/template"
	"net/http"
)

// Server is a server that presents the dashboard.
type Server struct {
	http.Server
	mux     *http.ServeMux
	conf    *config.Configuration
	counter *hitcounter.HitCounter
	t       *template.Template
}

// New returns an initialized instance of *Server.
func New(conf *config.Configuration, counter *hitcounter.HitCounter) *Server {
	result := new(Server)
	result.conf = conf
	result.counter = counter
	result.mux = http.NewServeMux()
	result.t = template.Must(template.New("core").Parse(coreHtml))
	result.Server = http.Server{
		Addr:    conf.DashboardAddress,
		Handler: result.mux,
	}

	result.setupRoutes()

	return result
}

// ListenAndServe starts the server in a new goroutine (non-blocking).
func (s *Server) ListenAndServe() {
	go s.Server.ListenAndServe()
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/history", s.HandleHistory)
	s.mux.HandleFunc("/", s.HandleAssets)
}

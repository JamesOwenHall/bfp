package dashboard

import (
	"github.com/JamesOwenHall/BruteForceProtection/core/config"
	"github.com/JamesOwenHall/BruteForceProtection/core/hitcounter"
	"html/template"
	"net/http"
)

type Server struct {
	http.Server
	mux     *http.ServeMux
	conf    *config.Configuration
	counter *hitcounter.HitCounter
	t       *template.Template
}

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

func (s *Server) ListenAndServe() {
	go s.Server.ListenAndServe()
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/history", s.HandleHistory)
	s.mux.HandleFunc("/", s.HandleAssets)
}

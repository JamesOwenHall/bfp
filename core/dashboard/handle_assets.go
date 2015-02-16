package dashboard

import (
	"github.com/JamesOwenHall/BruteForceProtection/core/config"
	"github.com/JamesOwenHall/BruteForceProtection/core/hitcounter"
	"html/template"
	"net/http"
)

func (s *Server) HandleAssets(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/core.css":
		http.ServeFile(w, r, "dashboard/static/core.css")
	case "/core.js":
		http.ServeFile(w, r, "dashboard/static/core.js")
	case "/":
		s.serveHomePage(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) serveHomePage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		ListenAddress string
		ListenType    string
		Version       string
		Directions    []hitcounter.Direction
	}{
		ListenAddress: s.conf.ListenAddress,
		ListenType:    s.conf.ListenType,
		Version:       config.Version,
		Directions:    s.conf.Directions,
	}

	t := template.Must(template.ParseFiles("dashboard/static/core.html"))
	t.Execute(w, data)
}

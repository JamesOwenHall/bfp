package dashboard

import (
	"github.com/JamesOwenHall/bfp/core/config"
	"github.com/JamesOwenHall/bfp/core/hitcounter"
	"net/http"
)

// HandleAssets serves the HTML, CSS and JS assets for the dashboard.
func (s *Server) HandleAssets(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/core.css":
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(coreCss))
	case "/core.js":
		w.Header().Set("Content-Type", "application/js")
		w.Write([]byte(coreJs))
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

	s.t.Execute(w, data)
}

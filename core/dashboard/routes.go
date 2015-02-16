package dashboard

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/history", s.HandleHistory)
	s.mux.HandleFunc("/", s.HandleAssets)
}

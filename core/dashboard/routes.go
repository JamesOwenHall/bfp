package dashboard

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/", s.HandleAssets)
}

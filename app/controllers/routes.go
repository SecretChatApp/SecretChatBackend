package controllers

func (s *Server) InitializeRoutes() {
	s.Router.HandleFunc("/", s.Home).Methods("GET")
}

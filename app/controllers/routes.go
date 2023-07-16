package controllers

func (s *Server) InitializeRoutes() {
	s.Router.HandleFunc("/", s.Home).Methods("GET")

	// Auth
	s.Router.HandleFunc("/register", s.Register).Methods("POST")
	s.Router.HandleFunc("/login", s.Login).Methods("POST")
	s.Router.HandleFunc("/logout", s.Logout).Methods("GET")
}

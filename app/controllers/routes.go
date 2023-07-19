package controllers

import "backend/app/middlewares"

func (s *Server) InitializeRoutes() {
	s.Router.HandleFunc("/", s.Home).Methods("GET")

	// Auth
	s.Router.HandleFunc("/register", s.Register).Methods("POST")
	s.Router.HandleFunc("/login", s.Login).Methods("POST")
	s.Router.HandleFunc("/logout", s.Logout).Methods("GET")

	api := s.Router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/chatrooms", s.GetChatrooms).Methods("GET")
	api.HandleFunc("/chatroom", s.CreateChatroom).Methods("POST")
	api.HandleFunc("/chatroom/{id}", s.GetRoomInformation).Methods("GET")
	api.HandleFunc("/edit/{id}", s.EditRoom).Methods("PUT")
	api.Use(middlewares.JWTMiddleware)

}

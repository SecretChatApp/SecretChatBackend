package controllers

import (
	"backend/app/middlewares"
	chatservices "backend/app/services/chat_services"
	"net/http"
)

func (s *Server) InitializeRoutes() {

	WsServer := chatservices.NewWebSocketServer()
	go WsServer.Run()

	s.Router.HandleFunc("/", s.Home).Methods("GET")

	// Auth
	s.Router.HandleFunc("/register", s.Register).Methods("POST")
	s.Router.HandleFunc("/login", s.Login).Methods("POST")
	s.Router.HandleFunc("/logout", s.Logout).Methods("GET")
	s.Router.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		s.ServeWs(WsServer, w, r)
	})

	api := s.Router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/chatrooms", s.GetChatrooms).Methods("GET")
	api.HandleFunc("/chatroom", s.CreateChatroom).Methods("POST")
	api.HandleFunc("/chatroom/{id}", s.GetRoomInformation).Methods("GET")
	api.HandleFunc("/edit/{id}", s.EditRoom).Methods("PUT")
	api.HandleFunc("/delete/{id}", s.DeleteRoom).Methods("DELETE")
	api.Use(middlewares.JWTMiddleware)

}

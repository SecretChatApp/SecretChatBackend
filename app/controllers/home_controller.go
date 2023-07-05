package controllers

import (
	"fmt"
	"net/http"
)

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello this is home")
}

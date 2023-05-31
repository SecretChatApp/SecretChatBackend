package controllers

import (
	"backend/app/helpers"
	"encoding/json"
	"net/http"
)


type UserInput struct {
  Name string `json:"name" validate:"require,gte=4"`
  Email string `json:"email" validate:"required,email,isunique=users-email"`
  Password string `json:"password" validate:"required,gte=4"`
  ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password" label:"Konfirmasi Password"`
}

type UserLogin struct {
  Email string `json:"email" validate:"required"`
  Password string `json:"password" validate:"required"`
}

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
  var userInput UserInput
  decoder := json.NewDecoder(r.Body)
  if err := decoder.Decode(&userInput); err != nil {
    response := map[string]string{
      "message": err.Error(),
    }

    helpers.ResponseJSON(w, http.StatusInternalServerError, response)
    return
  }

  defer r.Body.Close()
}

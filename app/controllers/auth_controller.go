package controllers

import (
	"backend/app/helpers"
	"backend/app/libraries"
	"backend/app/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	Name            string `json:"name" validate:"require,gte=4"`
	Email           string `json:"email" validate:"required,email,isunique=users-email"`
	Password        string `json:"password" validate:"required,gte=4"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password" label:"Konfirmasi Password"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required"`
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

	validate := libraries.NewValidation()
	errData := validate.Struct(userInput)
	if errData != nil {
		data := map[string]interface{}{
			"validation": errData,
			"user":       userInput,
		}

		helpers.ResponseJSON(w, http.StatusInternalServerError, data)
		return
	} else {
		hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
		userInput.Password = string(hashPassword)

		user := models.User{
			ID:       uuid.New().String(),
			Name:     userInput.Name,
			Email:    userInput.Email,
			Password: userInput.Password,
		}

		insertId, err := user.CreateUser(s.DB)
		var message string
		if err != nil {
			message = "Gagal Mendaftar"
			data := map[string]string{
				"message": message,
			}
			helpers.ResponseJSON(w, http.StatusInternalServerError, data)
			return
		} else {
			message = "Berhasil, Silahkan login"
			fmt.Println("User created with id: ", insertId)
			data := map[string]string{
				"message": message,
			}
			helpers.ResponseJSON(w, http.StatusOK, data)
		}
	}
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var userLogin UserLogin
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userLogin); err != nil {
		response := map[string]string{
			"message": err.Error(),
		}
		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	} else {
		// var user models.User{}

	}
}

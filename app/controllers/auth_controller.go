package controllers

import (
	"backend/app/config"
	"backend/app/helpers"
	"backend/app/libraries"
	"backend/app/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserInput struct {
	Name            string `json:"name" validate:"required,gte=4"`
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
		response := map[string]interface{}{
			"message": err.Error(),
		}

		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	defer r.Body.Close()

	validate := libraries.NewValidation()
	errData := validate.Struct(userLogin)
	if errData != nil {
		response := map[string]interface{}{
			"validation": errData,
		}

		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	} else {
		var user models.User
		err := user.GetUser(s.DB, userLogin.Email)
		var message string
		if err != nil {
			message = "User tidak ditemukan"
		} else {
			errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
			if errPassword != nil {
				message = "Password tidak valid"
			}
		}

		if message != "" {
			response := map[string]interface{}{
				"message": message,
			}
			helpers.ResponseJSON(w, http.StatusInternalServerError, response)
			return
		} else {
			expTime := time.Now().Add(time.Hour * 24)
			claims := &config.JWTClaim{
				Email: userLogin.Email,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "backend",
					ExpiresAt: jwt.NewNumericDate(expTime),
				},
			}

			tokenAlgo := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			token, err := tokenAlgo.SignedString(config.JWT_KEY)
			if err != nil {
				response := map[string]interface{}{
					"message": err.Error(),
				}

				helpers.ResponseJSON(w, http.StatusInternalServerError, response)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Path:     "/",
				Value:    token,
				HttpOnly: true,
				Expires:  expTime,
			})

			response := map[string]interface{}{
				"message": "Login Berhasil",
			}

			helpers.ResponseJSON(w, http.StatusOK, response)
			return
		}
	}
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	response := map[string]interface{}{
		"message": "Logout Berhasil",
	}

	helpers.ResponseJSON(w, http.StatusOK, response)
	return
}

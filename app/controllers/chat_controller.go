package controllers

import (
	"backend/app/helpers"
	"backend/app/libraries"
	"backend/app/models"
	chatservices "backend/app/services/chat_services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (s *Server) GetChatrooms(w http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value("userInfo").(string)

	var user models.User
	err := user.GetAllChatRoom(s.DB, userInfo)
	if err != nil {
		log.Fatal(err)
	}

	response := map[string]interface{}{
		"data": user.ChatRoom,
	}

	helpers.ResponseJSON(w, http.StatusOK, response)
	return
}

func (s *Server) CreateChatroom(w http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value("userInfo").(string)

	var inputChatroom chatservices.InputChatroom
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&inputChatroom); err != nil {
		response := map[string]interface{}{
			"message": err.Error(),
		}

		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	defer r.Body.Close()

	validate := libraries.NewValidation()
	errData := validate.Struct(inputChatroom)
	if errData != nil {
		response := map[string]interface{}{
			"message": errData,
		}

		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	var user models.User
	id, err := user.GetUserID(s.DB, userInfo)
	if err != nil {
		log.Fatal(err)
	}

	chatroom := models.ChatRoom{
		ID:      uuid.New().String(),
		UserID:  id,
		Title:   inputChatroom.Title,
		Subject: inputChatroom.Subject,
	}

	title, err := chatroom.CreateChatRoom(s.DB)

	var message string

	if err != nil {
		message = "Gagal membuat chatroom"
		response := map[string]interface{}{
			"message": message,
		}

		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	} else {
		message = "Berhasil membuat Chatroom dengan judul: " + title
		response := map[string]interface{}{
			"message": message,
		}

		helpers.ResponseJSON(w, http.StatusOK, response)
		return
	}
}

func (s *Server) EditRoom(w http.ResponseWriter, r *http.Request) {
	// name, ok := r.URL.Query()["name"]
	// if !ok || len(name[0]) < 1 {
	// 	log.Println("url param name is missing")
	// 	return
	// }

	params := mux.Vars(r)
	objID := params["id"]

	if len(objID) < 1 {
		log.Println("url param id is missing")
		return
	}

	var editRequest chatservices.EditRequest

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&editRequest); err != nil {
		response := map[string]string{
			"message": err.Error(),
		}

		helpers.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	defer r.Body.Close()

	var chatroom models.ChatRoom
	err1 := chatroom.GetRoomById(objID, s.DB)

	if err1 != nil {
		log.Println(err1)
		return
	}

	err := chatroom.UpdateRoom(editRequest.Title, editRequest.Subject, s.DB)
	if err != nil {
		log.Println(err)
		return
	}

	response := map[string]string{
		"message": "Berhasil update",
	}

	helpers.ResponseJSON(w, http.StatusOK, response)
	return
}

func (s *Server) GetRoomInformation(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	objID := param["id"]

	fmt.Println(objID)

	if len(objID) < 1 {
		log.Println("url param id is missing")
		return
	}

	var chatroom models.ChatRoom
	// chatroom.GetRoomInformation(name[0], s.DB)
	chatroom.GetRoomInformation(objID, s.DB)

	responseInformation := chatservices.ResponseRoomInformation{
		OwnerName: chatroom.User.Name,
		Title:     chatroom.Title,
		Subject:   chatroom.Subject,
		Messages:  chatroom.Message,
		CreatedAt: chatroom.CreatedAt,
	}

	response := map[string]interface{}{
		"data": responseInformation,
	}

	helpers.ResponseJSON(w, http.StatusOK, response)
	return
}

func (s *Server) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	name, ok := r.URL.Query()["name"]
	if !ok || len(name[0]) < 1 {
		log.Println("url param name is missing")
		return
	}

	var chatroom models.ChatRoom
	err := chatroom.RemoveChatRoom(name[0], s.DB)
	if err != nil {
		log.Print(err)
		return
	}

	response := map[string]string{
		"message": "Berhasil hapus room",
	}

	helpers.ResponseJSON(w, http.StatusOK, response)
	return
}

func (s *Server) ServeWs(wsServer *chatservices.WsServer, w http.ResponseWriter, r *http.Request) {
	upgrader := chatservices.Upgrader

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	_, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	name, ok := r.URL.Query()["name"]

	if !ok || len(name[0]) < 1 {
		log.Println("url param name is missing")
		return
	}

}

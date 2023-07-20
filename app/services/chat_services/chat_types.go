package chatservices

import (
	"backend/app/models"
	"time"

	"github.com/gorilla/websocket"
)

type WsServer struct {
	RegisterRoom   chan *Room
	UnregisterRoom chan *Room
	Register       chan *Client
	Unregister     chan *Client
	Rooms          map[*Room]bool
}

type Client struct {
	Name     string `json:"name"`
	Conn     *websocket.Conn
	WsServer *WsServer
	Send     chan []byte
	Rooms    *Room
}

type Message struct {
	Action    string    `json:"action"`
	Message   string    `json:"message"`
	Target    string    `json:"target"`
	Sender    string    `json:"sender"`
	CreatedAt time.Time `json:"created_at"`
}

type Room struct {
	Name       string
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

type InputChatroom struct {
	Title   string `json:"title" validate:"required"`
	Subject string `json:"subject" validate:"required"`
}

type EditRequest struct {
	Title   string
	Subject string
}

type ResponseRoomInformation struct {
	OwnerName string           `json:"owner_name"`
	Title     string           `json:"title"`
	Subject   string           `json:"subject"`
	Messages  []models.Message `json:"messages"`
	CreatedAt time.Time        `json:"created_at"`
}

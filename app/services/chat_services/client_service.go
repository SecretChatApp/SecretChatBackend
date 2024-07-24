package chatservices

import (
	"backend/app/models"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

const (
	// max wait time when writing message to peer
	writeWait = 10 * time.Second

	// max time till next pong from peer
	pongWait = 60 * time.Second

	// send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// maximum message size allowed from peer
	maxMessageSize
)

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func NewClient(conn *websocket.Conn, ws *WsServer, name string, room *Room) *Client {
	return &Client{
		Name:     name,
		Rooms:    room,
		Conn:     conn,
		WsServer: ws,
		Send:     make(chan []byte, 256),
	}
}

func (client *Client) GetName() string {
	return client.Name
}

func (client *Client) ReadPump(r *http.Request, db *gorm.DB) {
	defer func() {
		client.Disconnect()
	}()

	client.Conn.SetReadLimit(int64(maxMessageSize))
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, jsonMessage, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}
		client.HandleNewMessage(jsonMessage, r, db)
	}
}

func (client *Client) WritePump(db *gorm.DB) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			n := len(client.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) HandleNewMessage(jsonMessage []byte, r *http.Request, db *gorm.DB) {
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Println(err)
	}

	_, err := r.Cookie("access_token")
	if err != nil {
		message.Sender = "client"
	} else {
		message.Sender = "owner"
	}

	time := time.Now()
	message.CreatedAt = time

	var messageModel = models.Message{
		ID:         uuid.New().String(),
		ChatRoomID: message.Target,
		Text:       message.Text,
		Sender:     message.Sender,
		CreatedAt:  message.CreatedAt,
	}

	var err2 error
	err2 = messageModel.CreateMessage(db)
	if err2 != nil {
		log.Println(err2)
		return
	}

	switch message.Action {
	case SendMessageAction:
		roomName := message.Target
		if room := client.WsServer.FindRoomByName(roomName); room != nil {
			room.Broadcast <- &message
		}
	}
}

func (client *Client) Disconnect() {
	client.WsServer.Unregister <- client
	close(client.Send)
	client.Conn.Close()
}

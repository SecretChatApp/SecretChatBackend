package chatservices

import "github.com/gorilla/websocket"

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

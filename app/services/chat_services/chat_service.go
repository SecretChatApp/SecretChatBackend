package chatservices

import (
	"fmt"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func NewWebSocketServer() *WsServer {
	return &WsServer{
		RegisterRoom:   make(chan *Room),
		UnregisterRoom: make(chan *Room),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		Rooms:          make(map[*Room]bool),
	}
}

func (server *WsServer) CreateRoom(name string) *Room {
	var room *Room

	if res := server.CheckRoomByName(name); res {
		room = server.FindRoomByName(name)
	} else {
		room = NewRoom(name)
		server.Rooms[room] = true
	}

	go room.RunRoom()

	return room

}

func (server *WsServer) CheckRoomByName(name string) bool {
	for room := range server.Rooms {
		if room.GetRoomName() == name {
			return true
		}
	}

	return false
}

func (server *WsServer) FindRoomByName(name string) *Room {
	var foundRoom *Room
	for room := range server.Rooms {
		if room.GetRoomName() == name {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (server *WsServer) PrintAllRooms() {
	for room := range server.Rooms {
		fmt.Println(room.Name)
	}
}

func (server *WsServer) RegisteringRoom(room *Room) {
	if res := server.CheckRoomByName(room.Name); res {
		fmt.Println("Room already exists")
		return
	} else {
		server.Rooms[room] = true
		server.PrintAllRooms()
	}
}

func (server *WsServer) UnregisteringRoom(room *Room) {
	if res := server.CheckRoomByName(room.Name); res {
		delete(server.Rooms, room)
		server.PrintAllRooms()
	}
}

func (server *WsServer) Run() {
	for {
		select {
		case room := <-server.RegisterRoom:
			server.RegisteringRoom(room)
		case room := <-server.UnregisterRoom:
			server.UnregisteringRoom(room)
		}
	}
}

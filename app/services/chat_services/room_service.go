package chatservices

import "fmt"

func NewRoom(name string) *Room {
	return &Room{
		Name:       name,
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
	}
}

func (room *Room) RegisterClientInRoom(client *Client) {
	room.Clients[client] = true
	fmt.Println(client.Rooms.Name)
}

func (room *Room) PrintAllClients() {
	for client := range room.Clients {
		fmt.Println(&client.WsServer, &client.Conn)
	}
}

func (room *Room) UnregisterClientInRoom(client *Client) {
	if _, ok := room.Clients[client]; ok {
		delete(room.Clients, client)
	}
}

func (room *Room) BroadcastToClientsInRoom(message []byte) {
	for client := range room.Clients {
		client.Send <- message
	}
}

func (room *Room) GetRoomName() string {
	return room.Name
}

func (room *Room) RunRoom() {
	for {
		select {
		case client := <-room.Register:
			room.RegisterClientInRoom(client)
		case client := <-room.Unregister:
			room.UnregisterClientInRoom(client)
		case message := <-room.Broadcast:
			room.BroadcastToClientsInRoom(message.Encode())
		}
	}
}

package server

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[*websocket.Conn]bool

	Broadcast chan interface{}

	mu sync.Mutex // To protect the map from concurrent writes 
}


func NewHub() *Hub{
	return  &Hub{
		clients: make(map[*websocket.Conn]bool),
		Broadcast: make(chan interface{}),
	}
}


func (h *Hub) Run(){

		for {
			msg:= <-h.Broadcast
			h.mu.Lock()
			for client:= range h.clients {
				err:=client.WriteJSON(msg)
				if err!=nil{
					client.Close()
					delete(h.clients,client)
				}

			}
			h.mu.Unlock()
		}
}

func (h *Hub) Register(conn *websocket.Conn) {
	h.mu.Lock()
	h.clients[conn]=true
	h.mu.Unlock()
}
package websocket

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/polyk005/message/internal/api/repository"
	"github.com/polyk005/message/internal/api/service"
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type Client struct {
	ID   string
	Send chan Message
}

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	services   *service.Service
}

func NewHub(db *sqlx.DB) *Hub {
	return &Hub{
		services: &service.Service{
			Message: repository.NewMessageRepository(db),
		},
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Printf("Client registered: %s", client.ID)
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				log.Printf("Client unregistered: %s", client.ID)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешаем все источники (для разработки)
	},
}

// HandleWebSocket обрабатывает WebSocket-запросы.
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &Client{
		ID:   "unique-client-id", // Замените на реальный ID клиента
		Send: make(chan Message, 256),
	}

	h.Register <- client

	go h.writePump(client, conn)
	go h.readPump(client, conn)
}

func (h *Hub) writePump(client *Client, conn *websocket.Conn) {
	defer func() {
		conn.Close()
		h.Unregister <- client
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := conn.WriteJSON(message); err != nil {
				log.Println("WebSocket write error:", err)
				return
			}
		}
	}
}

func (h *Hub) readPump(client *Client, conn *websocket.Conn) {
	defer func() {
		conn.Close()
		h.Unregister <- client
	}()

	for {
		var message Message
		if err := conn.ReadJSON(&message); err != nil {
			log.Println("WebSocket read error:", err)
			break
		}
		h.Broadcast <- message
	}
}

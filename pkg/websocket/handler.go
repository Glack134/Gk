package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &Client{
		ID:   uuid.New().String(),
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
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second)) // Таймаут записи
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
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("WebSocket read error:", err)
			}
			break
		}

		// Обработка сообщения
		switch msg.Type {
		case "new_message":
			log.Println("Received new message:", msg.Payload)

			// Сохраняем сообщение в базе данных
			payload, ok := msg.Payload.(map[string]interface{})
			if !ok {
				log.Println("Invalid payload format")
				continue
			}

			chatID, ok := payload["chat_id"].(float64)
			if !ok {
				log.Println("Invalid chat_id format")
				continue
			}

			userID, ok := payload["user_id"].(float64)
			if !ok {
				log.Println("Invalid user_id format")
				continue
			}

			content, ok := payload["content"].(string)
			if !ok {
				log.Println("Invalid content format")
				continue
			}

			//сохранения сообщения
			messageID, err := h.services.Message.SendMessage(int(chatID), int(userID), content)
			if err != nil {
				log.Println("Failed to send message:", err)
				continue
			}

			// Отправляем сообщение всем участникам чата
			h.Broadcast <- Message{
				Type: "new_message",
				Payload: map[string]interface{}{
					"chat_id":    chatID,
					"message_id": messageID,
					"user_id":    userID,
					"content":    content,
				},
			}
		}
	}
}

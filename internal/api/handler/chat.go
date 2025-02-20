package handler

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/polyk005/message/pkg/websocket"
)

func (h *Handler) createChat(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input struct {
		Usernames []string `json:"usernames"`
		ChatName  string   `json:"chat_name"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var participantIDs []int
	for _, username := range input.Usernames {
		participantID, err := h.services.Chat.GetUserIDByUsername(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		participantIDs = append(participantIDs, participantID)
	}

	participantIDs = append(participantIDs, userID)

	//проверка чата создан ли он или нет
	sort.Ints(participantIDs)
	existingChatID, err := h.services.Chat.FindExistingChat(participantIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking chat existence"})
		return
	}

	if existingChatID != 0 {
		c.JSON(http.StatusOK, gin.H{"chat_id": existingChatID})
		return
	}

	chatID, err := h.services.Chat.CreateChat(input.ChatName, participantIDs...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.hub.Broadcast <- websocket.Message{
		Type:    "new_chat",
		Payload: map[string]interface{}{"chat_id": chatID, "user_ids": participantIDs},
	}

	c.JSON(http.StatusOK, gin.H{"chat_id": chatID})
}

func (h *Handler) getChatsForUser(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	chats, err := h.services.Chat.GetChatsForUser(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, chats)
}

func (h *Handler) addParticipant(c *gin.Context) {
	var input struct {
		ChatID        int `json:"chat_id"`
		ParticipantID int `json:"participant_id"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.Chat.AddParticipant(input.ChatID, input.ParticipantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "participant added"})
}

func (h *Handler) deleteChat(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input struct {
		ChatID int  `json:"chat_id"`
		ForAll bool `json:"for_all"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.ForAll {
		// Удаляем чат для всех участников
		err = h.services.Chat.DeleteChatForAll(input.ChatID)
	} else {
		// Удаляем чат только для текущего пользователя
		err = h.services.Chat.DeleteChatForUser(input.ChatID, userID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Отправляем уведомление через WebSocket
	h.hub.Broadcast <- websocket.Message{
		Type:    "delete_chat",
		Payload: map[string]interface{}{"chat_id": input.ChatID, "for_all": input.ForAll},
	}

	c.JSON(http.StatusOK, gin.H{"message": "chat deleted"})
}

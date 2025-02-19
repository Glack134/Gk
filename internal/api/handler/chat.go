package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createChat(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input struct {
		Username string `json:"username"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, существует ли пользователь
	userExists, err := h.services.Chat.UserExists(input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !userExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found. Please invite the user to register."})
		return
	}

	// Получаем ID пользователя, с которым создается чат
	participantID, err := h.services.Chat.GetUserIDByUsername(input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, существует ли уже чат между этими двумя пользователями
	chatID, err := h.services.Chat.ChatExistsBetweenUsers(userID, participantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if chatID != 0 {
		c.JSON(http.StatusOK, gin.H{"chat_id": chatID})
		return
	}

	// Если чат не существует, создаем новый
	chatID, err = h.services.Chat.CreateChat(userID, participantID, "New Chat")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

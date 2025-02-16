package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createChat(c *gin.Context) {
	var input struct {
		Name string `json:"name"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatID, err := h.services.Chat.CreateChat(input.Name)
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

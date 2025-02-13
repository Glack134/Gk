package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createChat(c *gin.Context) {
	var input struct {
		Name           string `json:"name"`
		ParticipantIDs []int  `json:"participant_ids"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatID, err := h.services.Chat.CreateChat(input.Name, input.ParticipantIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"chat_id": chatID})
}

func (h *Handler) getChat(c *gin.Context) {
	chatID := c.Param("id")

	chat, err := h.services.Chat.GetChat(chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chat)
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

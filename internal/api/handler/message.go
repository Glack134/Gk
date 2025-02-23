package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getMessages(c *gin.Context) {
	chatID := c.Param("chat_id")

	messages, err := h.services.Message.GetMessages(chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *Handler) sendMessage(c *gin.Context) {
	var input struct {
		ChatID  int    `json:"chat_id"`
		UserID  int    `json:"user_id"`
		Content string `json:"content"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messageID, err := h.services.Message.SendMessage(input.ChatID, input.UserID, input.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message_id": messageID})
}

func (h *Handler) editMessage(c *gin.Context) {
	var input struct {
		MessageID int    `json:"message_id"`
		Content   string `json:"content"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.Message.EditMessage(input.MessageID, input.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "message edited"})
}

func (h *Handler) deleteMessage(c *gin.Context) {
	messageID := c.Param("id")

	err := h.services.Message.DeleteMessage(messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "message deleted"})
}

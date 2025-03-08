package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getMessages(c *gin.Context) {
	chatIDStr := c.Param("chat_id")
	h.logger.Infof("Received chat_id: %s", chatIDStr) // Логируем значение chat_id

	if chatIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chat_id is required"})
		return
	}

	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		h.logger.Errorf("Invalid chat_id: %s", chatIDStr) // Логируем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid chat_id: must be an integer"})
		return
	}

	messages, err := h.services.Message.GetMessages(chatID)
	if err != nil {
		h.logger.Errorf("Failed to get messages: %v", err) // Логируем ошибку
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
	messageIDStr := c.Param("id")
	messageID, err := strconv.Atoi(messageIDStr) // Преобразуем строку в число
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message_id"})
		return
	}

	err = h.services.Message.DeleteMessage(messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "message deleted"})
}

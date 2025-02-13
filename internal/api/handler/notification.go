package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) sendNotification(c *gin.Context) {
	var input struct {
		UserID  int    `json:"user_id"`
		Message string `json:"message"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.Notification.SendNotification(input.UserID, input.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification sent"})
}

func (h *Handler) getNotifications(c *gin.Context) {
	userID := c.Param("id")

	notifications, err := h.services.Notification.GetNotifications(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createPayment(c *gin.Context) {
	var input struct {
		UserID  int     `json:"user_id"`
		Amount  float64 `json:"amount"`
		Purpose string  `json:"purpose"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentID, err := h.services.Payment.CreatePayment(input.UserID, input.Amount, input.Purpose)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment_id": paymentID})
}

func (h *Handler) getPaymentStatus(c *gin.Context) {
	paymentID := c.Param("id")

	status, err := h.services.Payment.GetPaymentStatus(paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}

package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
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

func (h *Handler) createSubscription(c *gin.Context) {
	var input struct {
		UserID int    `json:"user_id"`
		Plan   string `json:"plan"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscriptionID, err := h.services.Subscription.CreateSubscription(input.UserID, input.Plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subscription_id": subscriptionID})
}

func (h *Handler) getSubscription(c *gin.Context) {
	userID := c.Param("id")

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	subscription, err := h.services.Subscription.GetSubscription(userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

func (h *Handler) cancelSubscription(c *gin.Context) {
	subscriptionID := c.Param("id")

	subscriptionIDInt, err := strconv.Atoi(subscriptionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription ID"})
		return
	}

	err = h.services.Subscription.CancelSubscription(subscriptionIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subscription cancelled"})
}

func (h *Handler) generateQR(c *gin.Context) {
	// Генерация QR-кода с ссылкой на оплату
	qr, err := qrcode.Encode("https://your-payment-link.com", qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/png", qr)
}

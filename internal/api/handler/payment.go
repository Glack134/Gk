package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
)

func (h *Handler) createPayment(c *gin.Context) {
	var input struct {
		UserID   int     `json:"user_id"`
		Amount   float64 `json:"amount"`
		Purpose  string  `json:"purpose"`
		Currency string  `json:"currency"`
	}

	if err := c.BindJSON(&input); err != nil {
		h.logger.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создаем платеж в Stripe
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(input.Amount * 100)), // Сумма в центах
		Currency: stripe.String(strings.ToLower(input.Currency)),
	}
	pi, err := paymentintent.New(params)
	if err != nil {
		h.logger.Errorf("Failed to create payment intent: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем client_secret для подтверждения оплаты на фронтенде
	c.JSON(http.StatusOK, gin.H{
		"client_secret": pi.ClientSecret,
	})
}

func (h *Handler) generatePaymentLink(paymentID int, paymentMethod string, amount float64, currency string) (string, error) {
	switch paymentMethod {
	case "bank":
		return fmt.Sprintf("https://bank-payment-link.com/payments/%d?amount=%.2f&currency=%s", paymentID, amount, currency), nil
	case "paypal":
		return fmt.Sprintf("https://paypal.com/payment/%d?amount=%.2f&currency=%s", paymentID, amount, currency), nil
	case "stripe":
		return fmt.Sprintf("https://stripe.com/payment/%d?amount=%.2f&currency=%s", paymentID, amount, currency), nil
	default:
		return "", fmt.Errorf("unsupported payment method: %s", paymentMethod)
	}
}

func (h *Handler) handlePaymentCallback(c *gin.Context) {
	paymentIDStr := c.Query("payment_id")
	status := c.Query("status")

	// Проверка на пустые значения
	if paymentIDStr == "" || status == "" {
		h.logger.Warn("Missing payment_id or status in callback")
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment_id and status are required"})
		return
	}

	// Преобразуем paymentID из строки в число
	paymentID, err := strconv.Atoi(paymentIDStr)
	if err != nil {
		h.logger.Errorf("Invalid payment ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID"})
		return
	}

	// Обновляем статус платежа в базе данных
	err = h.services.Payment.UpdatePaymentStatus(paymentID, status)
	if err != nil {
		h.logger.Errorf("Failed to update payment status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Перенаправляем пользователя на страницу успешной оплаты
	c.Redirect(http.StatusFound, "/payment/success")
}

func (h *Handler) getPaymentStatus(c *gin.Context) {
	paymentID := c.Param("id")

	// Создаем контекст с таймаутом 5 секунд
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel() // Освобождаем ресурсы контекста

	// Передаем контекст в сервис
	status, err := h.services.Payment.GetPaymentStatus(ctx, paymentID)
	if err != nil {
		// Проверяем, была ли ошибка вызвана таймаутом
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "request timed out"})
			return
		}
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
	paymentIDStr := c.Param("id")
	paymentMethod := c.Query("method")

	paymentID, err := strconv.Atoi(paymentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment ID"})
		return
	}

	// Получаем данные платежа
	paymentDetails, err := h.services.Payment.GetPaymentDetails(paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Передаем amount и currency в generatePaymentLink
	paymentLink, err := h.generatePaymentLink(paymentID, paymentMethod, paymentDetails.Amount, paymentDetails.Currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	qr, err := qrcode.Encode(paymentLink, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "image/png", qr)
}

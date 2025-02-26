package handler

import (
	"image/png"
	"net/http"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/polyk005/message/internal/model"
)

func (h *Handler) signUp(c *gin.Context) {
	var input model.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Проверяем пароль и получаем пользователя
	user, err := h.services.Authorization.GetUser(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Проверяем, включена ли 2FA для пользователя
	isTwoFAEnabled, err := h.services.Authorization.IsTwoFAEnabled(user.Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to check 2FA status")
		return
	}

	if isTwoFAEnabled {
		// Если 2FA включена, возвращаем ответ с требованием ввести код 2FA
		c.JSON(http.StatusOK, gin.H{
			"message":      "2FA is enabled. Please provide the 2FA code.",
			"requires_2fa": true,
		})
		return
	}

	// Если 2FA не включена, генерируем токен и возвращаем его
	token, err := h.services.Authorization.GenerateToken(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":        token,
		"requires_2fa": false,
		"message":      "Login successful",
	})
}

func (h *Handler) verifyTwoFALogin(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Получаем пользователя по email
	user, err := h.services.Authorization.GetUser(input.Email, "")
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid email")
		return
	}

	// Проверяем код 2FA
	valid, err := h.services.Authorization.VerifyTwoFACode(user.Id, input.Code)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !valid {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid 2FA code")
		return
	}

	// Генерируем токен и возвращаем его
	token, err := h.services.Authorization.GenerateToken(input.Email, "")
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": "Login successful with 2FA",
	})
}

type requestResetPassword struct {
	Email string `json:"email" binding:"required"`
}

func (h *Handler) requestPasswordReset(c *gin.Context) {
	var input requestResetPassword

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err := h.services.SendPassword.CreateResetToken(input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Check to your email",
	})
}

func (h *Handler) ResetPasswordHandler(c *gin.Context) {
	if c.GetString("passwordResetDone") != "" {
		c.Redirect(http.StatusFound, "/main")
		return
	}

	token := c.Query("token")
	if token == "" {
		newErrorResponse(c, http.StatusBadRequest, "Token is required")
		return
	}

	// Проверяем токен
	err := h.services.Authorization.CheckToken(token)
	if err != nil {
		if err.Error() == "token has already been used" {
			c.Redirect(http.StatusFound, "/main")
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, "Internal server error")
		return
	}

	c.Set("passwordResetDone", true)

	c.HTML(http.StatusOK, "reset_password.html", gin.H{
		"token": token,
	})
}

func (h *Handler) UpdatePasswordHandler(c *gin.Context) {
	var input struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Пробуем привязать JSON-данные к структуре input
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Здесь вы можете проверить токен на валидность и срок действия
	// Если токен валиден, обновите пароль
	err := h.services.Authorization.UpdatePasswordUserToken(input.Token, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Возвращаем успешный ответ в формате JSON
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Пароль успешно обновлен",
	})
}

// 2fa
func (h *Handler) EnableTwoFA(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Генерируем новый секрет TOTP и получаем URL
	url, err := h.services.Authorization.EnableTwoFA(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Генерируем QR-код
	qrCode, err := qr.Encode(url, qr.L, qr.Auto)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to generate QR code")
		return
	}

	// Сохраняем QR-код в формате PNG
	qrCode, err = barcode.Scale(qrCode, 200, 200) // Масштабируем QR-код
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to scale QR code")
		return
	}

	// Устанавливаем заголовки для ответа
	c.Header("Content-Type", "image/png")

	// Отправляем QR-код в ответе
	if err := png.Encode(c.Writer, qrCode); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "failed to send QR code")
		return
	}

	// Формируем ответ с URL для ручного ввода кода
	response := map[string]interface{}{
		"manual_code": "",  // Здесь вы можете указать код, если он доступен
		"qr_code_url": url, // URL для QR-кода
	}

	// Отправляем JSON-ответ с кодом и URL для QR-кода
	c.JSON(http.StatusOK, response)
}

func (h *Handler) VerifyTwoFA(c *gin.Context) {
	var input struct {
		Code string `json:"code" binding:"required"`
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	valid, err := h.services.Authorization.VerifyTwoFACode(userId, input.Code)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if !valid {
		newErrorResponse(c, http.StatusUnauthorized, "invalid 2FA code")
		return
	}

	// Успешная проверка кода 2FA
	c.JSON(http.StatusOK, gin.H{"message": "2FA verification successful"})
}

func (h *Handler) DisableTwoFA(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	err = h.services.Authorization.DisableTwoFA(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Two-Factor Authentication disabled"})
}

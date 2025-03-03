package handler

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/polyk005/message/internal/model"
)

func validateKey(key []byte) error {
	switch len(key) {
	case 16, 24, 32:
		return nil
	default:
		return fmt.Errorf("invalid key size: %d (must be 16, 24, or 32 bytes)", len(key))
	}
}

func encrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func decrypt(ciphertext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	decodedCiphertext, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(decodedCiphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := decodedCiphertext[:aes.BlockSize]
	decodedCiphertext = decodedCiphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decodedCiphertext, decodedCiphertext)

	return string(decodedCiphertext), nil
}

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
		log.Printf("Ошибка привязки JSON: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "Неверный формат входных данных")
		return
	}

	log.Printf("Получены данные: %+v", input)

	user, err := h.services.Authorization.GetUser(input.Email, input.Password, true)
	if err != nil {
		log.Printf("Ошибка получения пользователя: %v", err)
		newErrorResponse(c, http.StatusUnauthorized, "Неверный email или пароль")
		return
	}

	log.Printf("Найден пользователь: %+v", user)

	isTwoFAEnabled, err := h.services.Authorization.IsTwoFAEnabled(user.Id)
	if err != nil {
		log.Printf("Ошибка проверки статуса 2FA: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, "Ошибка проверки статуса 2FA")
		return
	}

	if isTwoFAEnabled {
		key := []byte(os.Getenv("ENCRYPTION_KEY"))
		if err := validateKey(key); err != nil {
			log.Printf("Неверный ключ шифрования: %v", err)
			newErrorResponse(c, http.StatusInternalServerError, "Неверный ключ шифрования")
			return
		}

		encryptedUserID, err := encrypt(fmt.Sprintf("%d", user.Id), key)
		if err != nil {
			log.Printf("Ошибка шифрования ID пользователя: %v", err)
			newErrorResponse(c, http.StatusInternalServerError, "Ошибка шифрования ID пользователя")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      "2FA включен. Пожалуйста, введите код 2FA.",
			"requires_2fa": true,
			"user_id":      encryptedUserID,
		})
		return
	}

	accessToken, err := h.services.Authorization.GenerateAccessToken(user.Id)
	if err != nil {
		log.Printf("Failed to generate access token: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	refreshToken, err := h.services.Authorization.GenerateRefreshToken(user.Id)
	if err != nil {
		log.Printf("Failed to generate refresh token: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	domain := os.Getenv("COOKIE_DOMAIN")
	c.SetCookie("auth_token", accessToken, 3600, "/", domain, false, true)
	c.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", domain, false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"requires_2fa":  false,
		"message":       "Login successful",
	})
}

func (h *Handler) verifyTwoFALogin(c *gin.Context) {
	var input struct {
		EncryptedUserID string `json:"user_id" binding:"required"` // Изменено на string
		Code            string `json:"code" binding:"required"`
	}

	if err := c.BindJSON(&input); err != nil {
		log.Printf("Ошибка привязки JSON: %v", err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Дешифруем userId
	key := []byte(os.Getenv("ENCRYPTION_KEY"))
	decryptedUserID, err := decrypt(input.EncryptedUserID, key) // Теперь это строка
	if err != nil {
		log.Printf("Ошибка дешифрования: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "Failed to decrypt user ID")
		return
	}

	// Преобразуем расшифрованный ID в int
	userID, err := strconv.Atoi(decryptedUserID)
	if err != nil {
		log.Printf("Ошибка преобразования UserID: %v", err)
		newErrorResponse(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	// Проверяем код 2FA
	valid, err := h.services.Authorization.VerifyTwoFACode(userID, input.Code)
	if err != nil || !valid {
		log.Printf("Неверный код 2FA: %v", err)
		newErrorResponse(c, http.StatusUnauthorized, "Invalid 2FA code")
		return
	}

	// Генерируем токены
	accessToken, err := h.services.Authorization.GenerateAccessToken(userID)
	if err != nil {
		log.Printf("Ошибка генерации access token: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	refreshToken, err := h.services.Authorization.GenerateRefreshToken(userID)
	if err != nil {
		log.Printf("Ошибка генерации refresh token: %v", err)
		newErrorResponse(c, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	// Устанавливаем токены в куки
	domain := os.Getenv("COOKIE_DOMAIN")
	c.SetCookie("auth_token", accessToken, 3600, "/", domain, false, true)
	c.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", domain, false, true)

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"message":       "Login successful with 2FA",
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
		log.Printf("Failed to get user ID from token: %v", err)
		newErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	log.Printf("User ID retrieved from token: %d", userId)

	isTwoFAEnabled, err := h.services.Authorization.IsTwoFAEnabled(userId)
	if err != nil {
		log.Printf("Error checking 2FA status for user ID %d: %v", userId, err)
		newErrorResponse(c, http.StatusInternalServerError, "failed to check 2FA status")
		return
	}

	log.Printf("2FA status for user ID %d: %v", userId, isTwoFAEnabled)

	if isTwoFAEnabled {
		newErrorResponse(c, http.StatusBadRequest, "2FA is already enabled for this user")
		return
	}

	// Логика для включения 2FA...
	url, err := h.services.Authorization.EnableTwoFA(userId)
	if err != nil {
		log.Printf("Error enabling 2FA for user ID %d: %v", userId, err)
		newErrorResponse(c, http.StatusInternalServerError, "failed to enable 2FA")
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

	response := map[string]interface{}{
		"manual_code": "",  // Здесь вы можете указать код, если он доступен
		"qr_code_url": url, // URL для QR-кода
	}

	// Отправляем JSON-ответ с кодом и URL для QR-кода
	c.JSON(http.StatusOK, response)

	c.JSON(http.StatusOK, gin.H{
		"message": "2FA enabled successfully",
		"url":     url,
	})
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

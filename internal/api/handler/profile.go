package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/polyk005/message/internal/model"
)

const (
	signingKey = "qjvkvnsjdnj2njn29njv**@9un19@!33"
)

func (h *Handler) GetProfile(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.services.User.GetUserProfile(userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

func someHandler(c *gin.Context) {
	userId, err := getCurrentUserId(c, signingKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": userId})
}

func getCurrentUserId(c *gin.Context, signingKey string) (int, error) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return 0, fmt.Errorf("token not found")
	}

	// Удаляем "Bearer " из строки токена
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Далее ваш код для парсинга токена
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := int(claims["user_id"].(float64)) // Убедитесь, что ключ "user_id" соответствует вашему токену
		return userId, nil
	}

	fmt.Println("Invalid token:", err)
	return 0, err
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var input model.User_update
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUserId, err := getCurrentUserId(c, signingKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Проверка, что текущий пользователь обновляет свой профиль
	if input.Id != currentUserId {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own profile"})
		return
	}

	// Обновление профиля
	err = h.services.User.UpdateUserProfile(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func (h *Handler) RequestPasswordReset(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Здесь вы можете реализовать логику отправки кода на электронную почту
	// Например, с использованием SMTP или стороннего сервиса

	c.JSON(http.StatusOK, gin.H{"message": "Reset code sent to your email"})
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var input struct {
		Code     string `json:"code" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Здесь вы должны проверить код восстановления и затем обновить пароль
	// Например, с использованием базы данных для хранения временных кодов

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

func (h *Handler) UpdateEmail(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input struct {
		NewEmail string `json:"new_email" binding:"required,email"`
		Code     string `json:"code" binding:"required"`
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Проверка кода восстановления
	isValid, err := h.services.User.ValidateResetCode(input.Code) // Предполагается, что у вас есть такая функция
	if err != nil || !isValid {
		newErrorResponse(c, http.StatusBadRequest, "Invalid reset code")
		return
	}

	// Обновление электронной почты в базе данных
	if err := h.services.User.UpdateUserEmail(userID, input.NewEmail); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email updated successfully"})
}

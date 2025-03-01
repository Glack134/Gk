package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) refreshToken(c *gin.Context) {
	// Получаем refresh token из куки
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "Refresh token not found")
		return
	}

	// Проверяем refresh token
	userId, err := h.services.Authorization.ValidateRefreshToken(refreshToken)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	// Генерируем новый access token
	accessToken, err := h.services.Authorization.GenerateAccessToken(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Устанавливаем новый access token в куку
	c.SetCookie("auth_token", accessToken, 3600, "/", "localhost", false, true)          // Access token
	c.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", "localhost", false, true) // Refresh token

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"message":       "Token refreshed successfully",
	})
}

func (h *Handler) logOut(c *gin.Context) {
	// Удаляем access token
	c.SetCookie("auth_token", "", -1, "/", "localhost", false, true)

	// Удаляем refresh token
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

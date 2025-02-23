package handler

import (
	"fmt"
	"log"
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

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := int(claims["user_id"].(float64))
		return userId, nil
	}

	fmt.Println("Invalid token:", err)
	return 0, err
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var input model.User_update
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUserId, err := getCurrentUserId(c, signingKey)
	if err != nil {
		log.Printf("Failed to get user ID from token: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	input.Id = currentUserId

	err = h.services.User.UpdateUserProfile(&input)
	if err != nil {
		log.Printf("Failed to update profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

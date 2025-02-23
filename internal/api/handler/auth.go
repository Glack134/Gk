package handler

import (
	"net/http"

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

	token, err := h.services.Authorization.GenerateToken(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
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

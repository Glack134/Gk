package handler

import (
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/beefsack/go-rate"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

var (
	limiter   = rate.New(10, time.Second)
	mu        sync.Mutex
	bannedIPs = make(map[string]bool)
)

func (h *Handler) rateLimitMiddleware(c *gin.Context) {
	ip := c.ClientIP()

	//Исключаем localhost из rate limiting
	if ip == "::1" || ip == "127.0.0.1" {
		c.Next()
		return
	}

	mu.Lock()
	if bannedIPs[ip] {
		mu.Unlock()
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Ваш IP заблокирован"})
		return
	}
	mu.Unlock()

	go func() {
		for {
			time.Sleep(5 * time.Minute)
			mu.Lock()
			bannedIPs = make(map[string]bool)
			mu.Unlock()
		}
	}()

	ok, _ := limiter.Try()
	if !ok {
		mu.Lock()
		bannedIPs[ip] = true
		mu.Unlock()
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Слишком много запросов. Ваш IP заблокирован"})
		return
	}

	c.Next()
}

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found in context")
		return 0, errors.New("user id not found in context")
	}

	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id is of invalid type")
		return 0, errors.New("user id is of invalid type")
	}

	if idInt == 0 {
		newErrorResponse(c, http.StatusInternalServerError, "user id is 0")
		return 0, errors.New("user id is 0")
	}

	return idInt, nil
}

func (h *Handler) AuthMiddleware(c *gin.Context) {
	// Получаем токен из заголовка Authorization
	header := c.GetHeader("Authorization")
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	// Проверяем формат заголовка (Bearer <token>)
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	// Парсим токен
	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Устанавливаем user_id в контекст
	c.Set("user_id", userId)
	c.Next()
}

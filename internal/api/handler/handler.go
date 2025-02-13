package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/polyk005/message/internal/api/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	api := router.Group("/login")
	{
		api.POST("/sign-up", h.signUp)
		api.POST("/sign-in", h.signIn)
	}

	return router
}

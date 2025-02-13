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

	auth := router.Group("/auth")
	{
		auth.POST("/send-code", h.sendCode)
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	chat := router.Group("/chat")
	{
		chat.POST("/create", h.createChat)
		chat.GET("/:id", h.getChat)
		chat.POST("/add-participant", h.addParticipant)
	}

	message := router.Group("/message")
	{
		message.POST("/send", h.sendMessage)
		message.PUT("/edit", h.editMessage)
		message.DELETE("/:id", h.deleteMessage)
	}

	payment := router.Group("/payment")
	{
		payment.POST("/create", h.createPayment)
		payment.GET("/:id/status", h.getPaymentStatus)
	}

	notification := router.Group("/notification")
	{
		notification.POST("/send", h.sendNotification)
		notification.GET("/:id", h.getNotifications)
	}
	return router
}

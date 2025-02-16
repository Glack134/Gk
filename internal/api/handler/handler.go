package handler

import (
	"net/http"

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
	router := gin.Default()

	// Статические файлы
	router.Static("/static", "./frontend/static")

	// Маршруты для HTML страниц
	router.LoadHTMLGlob("frontend/*.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/login.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	router.GET("/signup.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", nil)
	})

	router.GET("/chat.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", nil)
	})

	router.GET("/profile.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "profile.html", nil)
	})

	// Группа для аутентификации
	auth := router.Group("/auth")
	{
		auth.POST("/send-code", h.sendCode)
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}

	chat := router.Group("/chat")
	chat.Use(h.AuthMiddleware())
	{
		chat.POST("/create", h.createChat)
		chat.GET("/:id", h.getChat)
		chat.POST("/add-participant", h.addParticipant)
	}

	message := router.Group("/message")
	message.Use(h.AuthMiddleware())
	{
		message.POST("/send", h.sendMessage)
		message.PUT("/edit", h.editMessage)
		message.DELETE("/:id", h.deleteMessage)
	}

	payment := router.Group("/payment")
	payment.Use(h.AuthMiddleware())
	{
		payment.POST("/create", h.createPayment)
		payment.GET("/:id/status", h.getPaymentStatus)
	}

	subscription := router.Group("/subscription")
	subscription.Use(h.AuthMiddleware())
	{
		subscription.POST("/create", h.createSubscription)
		subscription.GET("/:id", h.getSubscription)
		subscription.DELETE("/:id", h.cancelSubscription)
	}

	notification := router.Group("/notification")
	notification.Use(h.AuthMiddleware())
	{
		notification.POST("/send", h.sendNotification)
		notification.GET("/:id", h.getNotifications)
	}

	return router
}

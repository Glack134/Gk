package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/polyk005/message/internal/api/service"
	"github.com/polyk005/message/pkg/websocket"
)

type Handler struct {
	services *service.Service
	hub      *websocket.Hub
}

func NewHandler(services *service.Service, db *sqlx.DB) *Handler {
	hub := websocket.NewHub(db)
	go hub.Run()

	return &Handler{
		services: services,
		hub:      hub,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(h.rateLimitMiddleware)

	// Статические файлы
	router.Static("/static", "./frontend/static")

	// Маршруты для HTML страниц
	router.LoadHTMLGlob("frontend/*.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/ws", func(c *gin.Context) {
		h.hub.HandleWebSocket(c.Writer, c.Request)
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
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/reset_password", h.requestPasswordReset)
	}

	profile := router.Group("/profile")
	profile.Use(h.AuthMiddleware)
	{
		profile.GET("/", h.GetProfile)
		profile.PUT("/update", h.UpdateProfile)
		profile.POST("/enable", h.EnableTwoFA)
		profile.POST("/verify", h.VerifyTwoFA)
		profile.POST("/disable", h.DisableTwoFA)
	}

	chat := router.Group("/chat")
	chat.Use(h.AuthMiddleware)
	{
		chat.POST("/create", h.createChat)
		chat.POST("/chats", h.getChatsForUser)
		chat.POST("/add-participant", h.addParticipant)
		chat.DELETE("/delete", h.deleteChat)
	}

	message := router.Group("/message")
	message.Use(h.AuthMiddleware)
	{
		message.GET("/chat/:chat_id/messages", h.getMessages)
		message.POST("/send", h.sendMessage)
		message.PUT("/edit", h.editMessage)
		message.DELETE("/:id", h.deleteMessage)
	}

	payment := router.Group("/payment")
	payment.Use(h.AuthMiddleware)
	{
		payment.POST("/create", h.createPayment)
		payment.GET("/:id/status", h.getPaymentStatus)
	}

	subscription := router.Group("/subscription")
	subscription.Use(h.AuthMiddleware)
	{
		subscription.POST("/create", h.createSubscription)
		subscription.GET("/:id", h.getSubscription)
		subscription.DELETE("/:id", h.cancelSubscription)
	}

	notification := router.Group("/notification")
	notification.Use(h.AuthMiddleware)
	{
		notification.POST("/send", h.sendNotification)
		notification.GET("/:id", h.getNotifications)
	}

	return router
}

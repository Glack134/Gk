package handler

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/polyk005/message/internal/api/service"
	"github.com/polyk005/message/pkg/websocket"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/paymentintent"
)

type Handler struct {
	services *service.Service
	hub      *websocket.Hub
	logger   *logrus.Logger
}

func NewHandler(services *service.Service, db *sqlx.DB, logger *logrus.Logger) *Handler {
	hub := websocket.NewHub(db)
	go hub.Run()

	return &Handler{
		services: services,
		hub:      hub,
		logger:   logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost:8080"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders: []string{"Content-Length"},

		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(h.rateLimitMiddleware)

	// Группа для аутентификации
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/verify-2fa", h.verifyTwoFALogin)
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
		payment.GET("/callback", h.handlePaymentCallback)
		payment.GET("/:id/qr", h.generateQR)
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

	// Статические файлы
	router.Static("/static", "./frontend/static")

	router.LoadHTMLGlob("/home/midiy/file_programming/message/frontend/*.html")

	router.GET("/login2fa.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login2fa.html", nil)
	})

	router.GET("/ws", func(c *gin.Context) {
		h.hub.HandleWebSocket(c.Writer, c.Request)
	})

	router.GET("/", func(c *gin.Context) {
		// Проверяем, аутентифицирован ли пользователь
		_, err := c.Cookie("auth_token")
		if err != nil {
			c.Redirect(http.StatusFound, "login.html") // Перенаправляем на страницу входа
			return
		}

		c.Redirect(http.StatusFound, "chat.html") // Перенаправляем на чат, если пользователь аутентифицирован
	})

	router.GET("/login.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	router.GET("/signup.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", nil)
	})

	router.GET("/stripe", func(c *gin.Context) {
		// Создаем PaymentIntent
		params := &stripe.PaymentIntentParams{
			Amount:   stripe.Int64(100), // Сумма в центах (например, 1.00 USD)
			Currency: stripe.String("usd"),
		}
		pi, err := paymentintent.New(params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать PaymentIntent"})
			return
		}

		// Передаем client_secret в шаблон
		c.HTML(http.StatusOK, "stripe.html", gin.H{
			"client_secret": pi.ClientSecret, // Используем реальный client_secret
		})
	})

	// Защищенные маршруты
	router.GET("/chat.html", h.AuthMiddleware, func(c *gin.Context) {
		c.HTML(http.StatusOK, "chat.html", nil)
	})

	router.GET("/profile.html", h.AuthMiddleware, func(c *gin.Context) {
		c.HTML(http.StatusOK, "profile.html", nil)
	})

	router.POST("/auth/logout", func(c *gin.Context) {
		c.SetCookie("auth_token", "", -1, "/", "localhost", false, true)    // Удаляем access token
		c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true) // Удаляем refresh token
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	})

	return router
}

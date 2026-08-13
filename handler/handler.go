package handler

import (
	ws_hub "web-socket-postgres-auth-messager/handler/ws-hub"
	"web-socket-postgres-auth-messager/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
	hub      *ws_hub.Hub
}

func NewHandler(services *service.Service) *Handler {
	hub := ws_hub.NewHub()

	go hub.Run()

	return &Handler{services: services, hub: hub}
}

func (handler *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(setupHeaders())

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", handler.signUp)
		auth.POST("/sign-in", handler.signIn)
	}

	router.GET("/ws", handler.ws)

	return router
}

func setupHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Разрешаем передачу полного URL (отключаем strict-origin-when-cross-origin)
		c.Header("Referrer-Policy", "no-referrer-when-downgrade")

		// ОЧЕНЬ ВАЖНО для связки с Next.js: Настройка CORS
		// Укажите точный домен вашего Next.js (например, http://localhost:3000 в разработке)
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		// Если это предзапрос (Preflight) от браузера Next.js, сразу возвращаем 200
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		// Передаем управление следующим обработчикам (маршрутам)
		c.Next()
	}
}

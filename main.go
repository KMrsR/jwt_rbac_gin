package main

import (
	"jwt_rbac_gin/handlers"
	"jwt_rbac_gin/middleware"
	"jwt_rbac_gin/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Маршрут для логина (публичный)
	r.POST("/login", handlers.LoginHandler)

	// Защищённые маршруты
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	{
		// Доступно всем авторизованным
		protected.GET("/user", handlers.UserHandler)

		// Доступно только администраторам
		protected.GET("/admin", middleware.RoleMiddleware(models.RoleAdmin), handlers.AdminHandler)
	}

	r.Run(":8080")
}

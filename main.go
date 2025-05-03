package main

import (
	"context"
	"fmt"
	"jwt_rbac_gin/handlers"
	"jwt_rbac_gin/middleware"
	"jwt_rbac_gin/models"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// r.Run(":8080")
	// Запускаем сервер в отдельной горутине
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		fmt.Println("🚀 Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server ListenAndServe error: %v\n", err)
		}
	}()

	// Ожидаем сигнал о завершении работы
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("\n🛑 Shutting down server...")

	// Контекст с таймаутом для завершения активных соединений
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Корректное завершение сервера
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v\n", err)
	}

	fmt.Println("✅ Server exited gracefully")
}

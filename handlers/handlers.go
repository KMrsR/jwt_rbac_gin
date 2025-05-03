package handlers

import (
	"jwt_rbac_gin/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// База пользователей (в реальности это будет БД)
var users = map[string]models.User{
	"john": {UserName: "john", Password: "pass123", Role: models.RoleUser},
	"anna": {UserName: "anna", Password: "pass456", Role: models.RoleAdmin},
}

// GenerateJWT генерирует токен с ролью пользователя
func GenerateJWT(username, role string) (string, error) {
	claims := &jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(models.JwtKey)
}

// LoginHandler — обработчик аутентификации
func LoginHandler(c *gin.Context) {
	var creds models.User
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, exists := users[creds.UserName]
	if !exists || user.Password != creds.Password {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := GenerateJWT(user.UserName, user.Role)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// AdminHandler — доступ только для админов
func AdminHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome admin!"})
}

// UserHandler — доступ для всех авторизованных
func UserHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello user!"})
}

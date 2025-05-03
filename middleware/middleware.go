package middleware

import (
	"jwt_rbac_gin/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// AuthMiddleware проверяет наличие JWT в заголовке
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Missing token"})
			return
		}

		// Bearer <token> → убираем префикс
		if len(tokenStr) > 7 && strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = tokenStr[7:]
		}

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return models.JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("claims", claims) // сохраняем claims для последующего использования
		c.Next()
	}
}

// RoleMiddleware проверяет, что у пользователя есть нужная роль
func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		role := (*claims.(*jwt.MapClaims))["role"].(string)

		if role != requiredRole {
			c.AbortWithStatusJSON(403, gin.H{"error": "Forbidden: insufficient permissions"})
			return
		}

		c.Next()
	}
}

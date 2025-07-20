package middleware

import (
	"net/http"
	"simple-gin-backend/internal/cache"
	"simple-gin-backend/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check if user ID exists in Redis
		exists, err := cache.RedisClient.Exists(cache.Ctx, claims.UserID.String()).Result()
		if err != nil || exists == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found or unauthorized"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)

		c.Next()
	}
}

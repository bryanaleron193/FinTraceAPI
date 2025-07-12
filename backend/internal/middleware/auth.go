package middleware

import (
	"net/http"
	"simple-gin-backend/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware is a middleware to protect routes
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			c.Abort()
			return
		}

		// Remove the "Bearer " prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the JWT
		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Store the user ID in the context for future use
		c.Set("user_id", claims.UserID)

		c.Next() // Continue to the next handler
	}
}
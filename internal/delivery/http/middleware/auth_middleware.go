package middleware

import (
	"net/http"
	"strings"

	"github.com/Hdeee1/go-register-login-profile/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Auth header is required"})
			return 
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization format"})
			return 
		}

		tokenString := parts[1]

		claims, err := jwt.ValidateToken(tokenString, secretKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return 
		}

		ctx.Set("user_id", claims.UserId)
		ctx.Next()
	}
}
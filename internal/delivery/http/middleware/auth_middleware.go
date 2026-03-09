package middleware

import (
	"net/http"
	"strings"

	"github.com/Hdeee1/go-register-login-profile/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(secretKey string, blacklist *jwt.TokenBlacklist) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Auth header is required"})
			return 
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return 
		}

		tokenString := parts[1]

		isBlacklist := blacklist.IsBlacklisted(tokenString)
		if isBlacklist == true {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has been invalidated"})
			return 
		}

		claims, err := jwt.ValidateToken(tokenString, secretKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return 
		}

		ctx.Set("user_id", claims.UserId)
		ctx.Next()
	}
}
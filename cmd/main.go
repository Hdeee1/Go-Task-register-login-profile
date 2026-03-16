package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Hdeee1/go-register-login-profile/internal/delivery/http"
	"github.com/Hdeee1/go-register-login-profile/internal/delivery/http/middleware"
	repository "github.com/Hdeee1/go-register-login-profile/internal/repository/mysql"
	"github.com/Hdeee1/go-register-login-profile/internal/usecase"
	"github.com/Hdeee1/go-register-login-profile/pkg/database"
	"github.com/Hdeee1/go-register-login-profile/pkg/jwt"
	"github.com/Hdeee1/go-register-login-profile/pkg/rabbitmq"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Failed to load env")
	}

	_, channel, err := rabbitmq.ConnectRabbitMQ(context.Background(), os.Getenv("RABBIT_URL"))
	if err != nil {
		log.Fatal("Failed to connect Rabbit MQ")
	}

	db, err := database.ConnectMySQL()
	if err != nil {
		log.Fatalf("Failed to connect database. Error: %s", err.Error())
	}

	repo, err := repository.NewUserRepository(db, logger)
	if err != nil {
		log.Fatal("Failed to create user repository")
	}

	useCase := usecase.NewUserUseCase(repo, logger, channel)
	blackList := jwt.NewTokenBlacklist()
	h := http.NewUserHandler(useCase, blackList, logger)

	rateLimiter := middleware.NewIPRateLimiter(1, 5)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	
	api := r.Group("/api")
	api.Use(middleware.RateLimiterMiddleware(rateLimiter))
	{
		api.POST("/user/register", h.Register)
		api.POST("/user/login", h.Login)
		api.POST("/auth/refresh", h.Refresh)
		api.POST("/auth/forgot-password", h.ForgotPassword)
		api.POST("/auth/reset-password", h.ResetPassword)

		auth := api.Group("/auth")
		auth.Use(middleware.AuthMiddleware(os.Getenv("JWT_ACCESS_SECRET"), blackList))
		{
			auth.GET("/profile", h.GetProfile)
			auth.PUT("/profile", h.UpdateProfile)
			auth.POST("/logout", h.Logout)
		}
	}

	fmt.Println("Server started at port :8080")
	r.Run(":8080")
}
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Hdeee1/go-register-login-profile/internal/delivery/http"
	"github.com/Hdeee1/go-register-login-profile/internal/delivery/http/middleware"
	repository "github.com/Hdeee1/go-register-login-profile/internal/repository/mysql"
	"github.com/Hdeee1/go-register-login-profile/internal/usecase"
	"github.com/Hdeee1/go-register-login-profile/pkg/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Failed to load env")
	}

	db, err := database.ConnectMySQL()
	if err != nil {
		log.Fatalf("Failed to connect database. Error: %s", err.Error())
	}

	repo, err := repository.NewUserRepository(db)
	if err != nil {
		log.Fatal("Failed to create user repository")
	}

	useCase := usecase.NewUserUsecase(repo)
	h := http.NewUserHandler(useCase)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"}, // adjust the frontend port
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
	
	api := r.Group("/api")
	{
		api.POST("/user/register", h.Register)
		api.POST("/user/login", h.Login)

		auth := api.Group("/auth")
		auth.Use(middleware.AuthMiddleware(os.Getenv("JWT_ACCESS_SECRET")))
		{
			auth.GET("/profile", h.GetProfile)
		}
	}

	fmt.Println("Server started at port :8080")
	r.Run(":8080")
}
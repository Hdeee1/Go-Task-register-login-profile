package http

import (
	"net/http"

	"github.com/Hdeee1/go-register-login-profile/internal/domain"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase domain.UserUsecase
}

func NewUserHandler(u domain.UserUsecase) *UserHandler {
	return &UserHandler{userUseCase: u}
}

func (h UserHandler) Register(ctx *gin.Context) error {
	var newUser domain.User

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "failed to bind json"})
	}

	ctx.JSON(http.StatusOK, newUser)
	return nil
}
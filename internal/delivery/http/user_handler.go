package http

import (
	"net/http"

	"github.com/Hdeee1/go-register-login-profile/internal/domain"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase domain.UserUsecase
}

type loginResponse struct {
	User 		 domain.User
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewUserHandler(u domain.UserUsecase) *UserHandler {
	return &UserHandler{userUseCase: u}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var newUser domain.User

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := h.userUseCase.Register(newUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var newUser domain.User

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	usr, accTkn, refTkn, err := h.userUseCase.Login(newUser)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	res := loginResponse{
		User: *usr,
		AccessToken: accTkn,
		RefreshToken: refTkn,
	}

	ctx.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *UserHandler) GetProfile(ctx *gin.Context) {
	value, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	userId := value.(int)

	user, err := h.userUseCase.GetProfile(userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}
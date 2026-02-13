package http

import (
	"net/http"

	"github.com/Hdeee1/go-register-login-profile/internal/domain"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase domain.UserUsecase
}

type registerResponse struct {
	Id			int `json:"user_id"`
	FullName	string `json:"full_name"`
	Username	string `json:"username"`
	Email		string `json:"email"`
}

type loginResponse struct {
	Username	 string `json:"username"`
	Email		 string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewUserHandler(u domain.UserUsecase) *UserHandler {
	return &UserHandler{userUseCase: u}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var newUser domain.RegisterRequest

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userUseCase.Register(newUser, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := registerResponse{
		Id: user.Id,
		FullName: user.FullName,
		Username: user.Username,
		Email: user.Email,
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": res})
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var newUser domain.LoginRequest

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usr, accTkn, refTkn, err := h.userUseCase.Login(newUser, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	res := loginResponse{
		Username: usr.Username,
		Email: usr.Email,
		AccessToken: accTkn,
		RefreshToken: refTkn,
	}

	ctx.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *UserHandler) GetProfile(ctx *gin.Context) {
	value, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userId := value.(int)

	user, err := h.userUseCase.GetProfile(userId, ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	res := gin.H{
		"data": gin.H{
			"id": user.Id,
			"full_name": user.FullName,
			"username": user.Username,
			"email": user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	}

	ctx.JSON(http.StatusOK, gin.H{"data": res})
}
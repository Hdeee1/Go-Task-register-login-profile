package http

import (
	"net/http"
	"os"
	"strings"

	"github.com/Hdeee1/go-register-login-profile/internal/domain"
	"github.com/Hdeee1/go-register-login-profile/pkg/jwt"
	"github.com/Hdeee1/go-register-login-profile/pkg/response"
	"github.com/Hdeee1/go-register-login-profile/pkg/validator"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
	tokenBlacklist *jwt.TokenBlacklist
}

type registerResponse struct {
	UserID		int	   `json:"user_id"`
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

func NewUserHandler(u domain.UserUseCase, b *jwt.TokenBlacklist) *UserHandler {
	return &UserHandler{
		userUseCase: u,
		tokenBlacklist: b,
	}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var newUser domain.RegisterRequest

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", validator.ParseValidatorError(err)))
		return
	}

	user, err := h.userUseCase.Register(newUser, ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", err.Error()))
		return
	}

	res := registerResponse{
		UserID: user.UserID,
		FullName: user.FullName,
		Username: user.Username,
		Email: user.Email,
	}
	ctx.JSON(http.StatusCreated, response.BuildSuccessResponse("CREATED", res))
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var newUser domain.LoginRequest

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", validator.ParseValidatorError(err)))
		return
	}

	usr, accTkn, refTkn, err := h.userUseCase.Login(newUser, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.BuildErrorResponse("UNAUTHORIZED", validator.ParseValidatorError(err)))
		return
	}

	res := loginResponse{
		Username: usr.Username,
		Email: usr.Email,
		AccessToken: accTkn,
		RefreshToken: refTkn,
	}
	ctx.JSON(http.StatusOK, response.BuildSuccessResponse("OK", res))
}

func (h *UserHandler) Logout(ctx *gin.Context) {
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

	claims, err := jwt.ValidateToken(tokenString, os.Getenv("JWT_ACCESS_SECRET"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return 
	}

	h.tokenBlacklist.AddTokenBlacklist(tokenString, claims.ExpiresAt.Time)
	ctx.JSON(http.StatusOK, response.BuildSuccessResponse("OK", gin.H{"message": "logged out"}))
}

func (h *UserHandler) Refresh(ctx *gin.Context) {
	var refresh domain.RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&refresh); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", validator.ParseValidatorError(err)))
		return
	}

	ref, err := h.userUseCase.Refresh(refresh, ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.BuildErrorResponse("UNAUTHORIZED", validator.ParseValidatorError(err)))
		return
	}
	ctx.JSON( http.StatusOK, response.BuildSuccessResponse("OK", gin.H{"access_token": ref}))
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
		ctx.JSON(http.StatusNotFound, response.BuildErrorResponse("NOT_FOUND", validator.ParseValidatorError(err)))
		return
	}

	res := gin.H{
			"id": user.UserID,
			"full_name": user.FullName,
			"username": user.Username,
			"email": user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
	}
	ctx.JSON(http.StatusOK, response.BuildSuccessResponse("OK", res))
}

func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	value, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userId := value.(int)

	var updateUser domain.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&updateUser); err != nil {
		ctx.JSON(http.StatusForbidden, response.BuildErrorResponse("BAD_REQUEST", validator.ParseValidatorError(err)))
		return
	}

	updatedUser, err := h.userUseCase.UpdateProfile(userId, updateUser, ctx)
	if err != nil {
		if err.Error() == "no fields to update" {
			ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", err.Error()))
			return
		}
		ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, response.BuildSuccessResponse("OK", &updatedUser))
}

func (h *UserHandler) ForgotPassword(ctx *gin.Context) {
	var forgotPass domain.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&forgotPass); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", validator.ParseValidatorError(err)))
		return
	}

	if err := h.userUseCase.ForgotPassword(forgotPass, ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, response.BuildSuccessResponse("The OTP code has been sent to your email", nil))
}

func (h *UserHandler) ResetPassword(ctx *gin.Context) {
	var reset domain.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&reset); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", err.Error()))
		return 
	}

	if err := h.userUseCase.ResetPassword(reset, ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, response.BuildSuccessResponse("The password has been changed", nil))
}
package http

import (
	"net/http"
	"os"
	"strings"

	"github.com/Hdeee1/go-register-login-profile/internal/delivery/http/dto"
	"github.com/Hdeee1/go-register-login-profile/internal/domain"
	"github.com/Hdeee1/go-register-login-profile/pkg/jwt"
	"github.com/Hdeee1/go-register-login-profile/pkg/response"
	"github.com/Hdeee1/go-register-login-profile/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	userUseCase domain.UserUseCase
	tokenBlacklist *jwt.TokenBlacklist
	logger *zap.Logger
}


func NewUserHandler(u domain.UserUseCase, b *jwt.TokenBlacklist, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userUseCase: u,
		tokenBlacklist: b,
		logger: logger,
	}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var newUser dto.RegisterRequest
	
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		h.logger.Warn("invalid request body on register", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", validator.ParseValidatorError(err)))
		return
	}
	
	user, err := h.userUseCase.Register(newUser, ctx)
	if err != nil {
		h.logger.Warn("register failed", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", err.Error()))
		return
	}
	
	res := dto.RegisterResponse{
		UserID: user.UserID,
		FullName: user.FullName,
		Username: user.Username,
		Email: user.Email,
	}
	h.logger.Info("user register successfully", zap.String("email", newUser.Email))
	ctx.JSON(http.StatusCreated, response.BuildSuccessResponse("CREATED", res))
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var newUser dto.LoginRequest
	
	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		h.logger.Warn("invalid request body on login", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", validator.ParseValidatorError(err)))
		return
	}
	
	usr, accTkn, refTkn, err := h.userUseCase.Login(newUser, ctx)
	if err != nil {
		h.logger.Warn("invalid login credentials", zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, response.BuildErrorResponse("UNAUTHORIZED", validator.ParseValidatorError(err)))
		return
	}

	res := dto.LoginResponse{
		Username: usr.Username,
		Email: usr.Email,
		AccessToken: accTkn,
		RefreshToken: refTkn,
	}
	h.logger.Info("user logged in successfully", zap.String("email", newUser.Identifier))
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
		h.logger.Warn("invalid token", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return 
	}
	h.logger.Info("user logged out", zap.Int("user_id", claims.UserId))

	h.tokenBlacklist.AddTokenBlacklist(tokenString, claims.ExpiresAt.Time)
	ctx.JSON(http.StatusOK, response.BuildSuccessResponse("OK", gin.H{"message": "logged out"}))
}

func (h *UserHandler) Refresh(ctx *gin.Context) {
	var refresh dto.RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&refresh); err != nil {
		h.logger.Warn("invalid request body on refresh token", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", validator.ParseValidatorError(err)))
		return
	}

	ref, err := h.userUseCase.Refresh(refresh, ctx)
	if err != nil {
		h.logger.Warn("failed to refresh token", zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, response.BuildErrorResponse("UNAUTHORIZED", validator.ParseValidatorError(err)))
		return
	}
	h.logger.Info("refresh token success")
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
		h.logger.Error("failed to Get Profile", zap.Error(err), zap.Int("user_id", userId))
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

	var updateUser dto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&updateUser); err != nil {
		h.logger.Warn("failed to Update Profile", zap.Error(err), zap.Int("user_id", userId))
		ctx.JSON(http.StatusForbidden, response.BuildErrorResponse("BAD_REQUEST", validator.ParseValidatorError(err)))
		return
	}

	updatedUser, err := h.userUseCase.UpdateProfile(userId, updateUser, ctx)
	if err != nil {
		h.logger.Warn("failed to update profile", zap.Error(err), zap.Int("user_id", userId))
		if err.Error() == "no fields to update" {
			ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", err.Error()))
			return
		}
		ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", err.Error()))
		return
	}
	h.logger.Info("update profile successfully", zap.Int("user_id", value.(int)))
	ctx.JSON(http.StatusOK, response.BuildSuccessResponse("OK", &updatedUser))
}

func (h *UserHandler) ForgotPassword(ctx *gin.Context) {
	var forgotPass dto.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&forgotPass); err != nil {
		h.logger.Warn("failed to request to forgot password", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", validator.ParseValidatorError(err)))
		return
	}

	if err := h.userUseCase.ForgotPassword(forgotPass, ctx); err != nil {
		h.logger.Warn("failed to request to forgot password", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", err.Error()))
		return
	}
	h.logger.Info("forgot password success", zap.String("email", forgotPass.Identifier))
	ctx.JSON(http.StatusOK, response.BuildSuccessResponse("The OTP code has been sent to your email", nil))
}

func (h *UserHandler) ResetPassword(ctx *gin.Context) {
	var reset dto.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&reset); err != nil {
		h.logger.Warn("failed to reset password", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", err.Error()))
		return 
	}

	if err := h.userUseCase.ResetPassword(reset, ctx); err != nil {
		h.logger.Warn("failed to reset password", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, response.BuildErrorResponse("BAD_REQUEST", err.Error()))
		return
	}
	h.logger.Info("reset password success", zap.String("email", reset.Identifier))
	ctx.JSON(http.StatusOK, response.BuildSuccessResponse("The password has been changed", nil))
}
package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shahariaz/gin-auth-service/internal/errs"
	"github.com/shahariaz/gin-auth-service/internal/model"
	"github.com/shahariaz/gin-auth-service/internal/service"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	service *service.AuthService
	log     *logrus.Logger
}

func NewAuthHandler(svc *service.AuthService, log *logrus.Logger) *AuthHandler {
	return &AuthHandler{service: svc, log: log}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required,min=3"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		errs.HandleValidationError(c, err, h.log)
		return
	}

	user := model.User{
		Username: input.Username,
		Email:    input.Email,
		RoleID:   1, // Default to 'user' role (ID 1)
	}
	if err := h.service.Register(&user, input.Password); err != nil {
		errs.HandleError(c, errs.NewAPIError(http.StatusBadRequest, "Registration failed", err), h.log)
		return
	}

	h.log.WithField("username", user.Username).Info("User registered")
	c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		errs.HandleValidationError(c, err, h.log)
		return
	}

	user, accessToken, refreshToken, err := h.service.Login(input.Email, input.Password)
	if err != nil {
		errs.HandleError(c, errs.NewAPIError(http.StatusUnauthorized, "Invalid credentials", err), h.log)
		return
	}

	h.log.WithField("username", user.Username).Info("User logged in")
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		errs.HandleValidationError(c, err, h.log)
		return
	}

	accessToken, err := h.service.RefreshToken(input.RefreshToken)
	if err != nil {
		errs.HandleError(c, errs.NewAPIError(http.StatusUnauthorized, "Invalid refresh token", err), h.log)
		return
	}

	h.log.Info("Token refreshed")
	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		errs.HandleValidationError(c, err, h.log)
		return
	}

	if err := h.service.Logout(input.RefreshToken); err != nil {
		errs.HandleError(c, errs.NewAPIError(http.StatusBadRequest, "Logout failed", err), h.log)
		return
	}

	h.log.Info("User logged out")
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

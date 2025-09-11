package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shahariaz/gin-auth-service/internal/errs"
	"github.com/shahariaz/gin-auth-service/internal/model"
	"github.com/shahariaz/gin-auth-service/internal/service"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	service *service.UserService
	log     *logrus.Logger
}

func NewUserHandler(svc *service.UserService, log *logrus.Logger) *UserHandler {
	return &UserHandler{service: svc, log: log}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	username, _ := c.Get("user")
	user, err := h.service.GetUserByUsername(username.(string))
	if err != nil {
		errs.HandleError(c, errs.NewAPIError(http.StatusNotFound, "User not found", err), h.log)
		return
	}
	h.log.WithField("username", user.Username).Info("Profile fetched")
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	username, _ := c.Get("user")
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		errs.HandleValidationError(c, err, h.log)
		return
	}

	user, err := h.service.UpdateUserProfile(username.(string), input.Email)
	if err != nil {
		errs.HandleError(c, errs.NewAPIError(http.StatusBadRequest, "Profile update failed", err), h.log)
		return
	}
	h.log.WithField("username", user.Username).Info("Profile updated")
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated", "user": user})
}

func (h *UserHandler) DeleteProfile(c *gin.Context) {
	username, _ := c.Get("user")
	if err := h.service.DeleteUser(username.(string)); err != nil {
		errs.HandleError(c, errs.NewAPIError(http.StatusBadRequest, "Profile deletion failed", err), h.log)
		return
	}
	h.log.WithField("username", username).Info("Profile deleted")
	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted"})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.service.ListUsers()
	if err != nil {
		errs.HandleError(c, errs.NewAPIError(http.StatusInternalServerError, "Failed to list users", err), h.log)
		return
	}
	h.log.Info("Fetched users list")
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		errs.HandleValidationError(c, err, h.log)
		return
	}
	if err := h.service.CreateUser(&user); err != nil {
		errs.HandleError(c, errs.NewAPIError(http.StatusBadRequest, "User creation failed", err), h.log)
		return
	}
	h.log.WithField("username", user.Username).Info("User created by admin")
	c.JSON(http.StatusCreated, gin.H{"message": "User created", "user": user})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errs.HandleError(c, errs.NewAPIError(http.StatusBadRequest, "Invalid user ID", err), h.log)
		return
	}
	var input struct {
		Username string `json:"username" binding:"required,min=3"`
		Email    string `json:"email" binding:"required,email"`
		RoleID   uint   `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		errs.HandleValidationError(c, err, h.log)
		return
	}

	user, err := h.service.UpdateUser(uint(id), input.Username, input.Email, input.RoleID)
	if err != nil {
		errs.HandleError(c, errs.NewAPIError(http.StatusBadRequest, "User update failed", err), h.log)
		return
	}
	h.log.WithField("username", user.Username).Info("User updated by admin")
	c.JSON(http.StatusOK, gin.H{"message": "User updated", "user": user})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errs.HandleError(c, errs.NewAPIError(http.StatusBadRequest, "Invalid user ID", err), h.log)
		return
	}
	if err := h.service.DeleteUserByID(uint(id)); err != nil {
		errs.HandleError(c, errs.NewAPIError(http.StatusBadRequest, "User deletion failed", err), h.log)
		return
	}
	h.log.WithField("id", id).Info("User deleted by admin")
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

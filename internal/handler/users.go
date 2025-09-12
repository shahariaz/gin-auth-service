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

// GetProfile godoc
// @Summary Get user profile
// @Description Get the authenticated user's profile information
// @Tags User Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "User profile retrieved successfully"
// @Failure 401 {object} map[string]string "Unauthorized - invalid or missing token"
// @Failure 404 {object} map[string]string "User not found"
// @Router /api/profile [get]
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

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the authenticated user's profile information
// @Tags User Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.UpdateProfileRequest true "Profile update request"
// @Success 200 {object} map[string]interface{} "Profile updated successfully"
// @Failure 400 {object} map[string]string "Bad request - validation error or update failed"
// @Failure 401 {object} map[string]string "Unauthorized - invalid or missing token"
// @Router /api/profile [put]
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

// DeleteProfile godoc
// @Summary Delete user profile
// @Description Delete the authenticated user's account
// @Tags User Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Profile deleted successfully"
// @Failure 400 {object} map[string]string "Bad request - profile deletion failed"
// @Failure 401 {object} map[string]string "Unauthorized - invalid or missing token"
// @Router /api/profile [delete]
func (h *UserHandler) DeleteProfile(c *gin.Context) {
	username, _ := c.Get("user")
	if err := h.service.DeleteUser(username.(string)); err != nil {
		errs.HandleError(c, errs.NewAPIError(http.StatusBadRequest, "Profile deletion failed", err), h.log)
		return
	}
	h.log.WithField("username", username).Info("Profile deleted")
	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted"})
}

// ListUsers godoc
// @Summary List all users (Admin only)
// @Description Get a list of all users in the system (requires admin role)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Users list retrieved successfully"
// @Failure 401 {object} map[string]string "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]string "Forbidden - admin role required"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/admin/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	users, err := h.service.ListUsers()
	if err != nil {
		errs.HandleError(c, errs.NewAPIError(http.StatusInternalServerError, "Failed to list users", err), h.log)
		return
	}
	h.log.Info("Fetched users list")
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// CreateUser godoc
// @Summary Create new user (Admin only)
// @Description Create a new user account (requires admin role)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateUserRequest true "User creation request"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]string "Bad request - validation error or user creation failed"
// @Failure 401 {object} map[string]string "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]string "Forbidden - admin role required"
// @Router /api/admin/users [post]
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

// UpdateUser godoc
// @Summary Update user (Admin only)
// @Description Update user information (requires admin role)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body model.UpdateUserRequest true "User update request"
// @Success 200 {object} map[string]interface{} "User updated successfully"
// @Failure 400 {object} map[string]string "Bad request - validation error or user update failed"
// @Failure 401 {object} map[string]string "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]string "Forbidden - admin role required"
// @Router /api/admin/users/{id} [put]
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

// DeleteUser godoc
// @Summary Delete user (Admin only)
// @Description Delete a user account (requires admin role)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string "User deleted successfully"
// @Failure 400 {object} map[string]string "Bad request - invalid user ID or deletion failed"
// @Failure 401 {object} map[string]string "Unauthorized - invalid or missing token"
// @Failure 403 {object} map[string]string "Forbidden - admin role required"
// @Router /api/admin/users/{id} [delete]
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

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/shahariaz/gin-auth-service/internal/config"
	"github.com/shahariaz/gin-auth-service/internal/database"
	"github.com/shahariaz/gin-auth-service/internal/handler"
	"github.com/shahariaz/gin-auth-service/internal/lib"
	"github.com/shahariaz/gin-auth-service/internal/middleware"
	"github.com/shahariaz/gin-auth-service/internal/service"
	"github.com/shahariaz/gin-auth-service/internal/validation"
	"github.com/sirupsen/logrus"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config, db *database.Database, log *logrus.Logger) {
	// Initialize dependencies/ // Switch to RedisTokenStore for prod

	tokenStore, err := lib.NewRedisTokenStore("localhost:6379")
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	validator := validation.NewValidator()
	userService := service.NewUserService(db, validator, log)
	authService := service.NewAuthService(db, validator, tokenStore, cfg.JWT_SECRET, log)
	userHandler := handler.NewUserHandler(userService, log)
	authHandler := handler.NewAuthHandler(authService, log)

	// Public routes
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/refresh", authHandler.RefreshToken)
	r.POST("/logout", authHandler.Logout)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.JWTAuthMiddleware(cfg.JWT_SECRET, tokenStore, log))
	{
		// User profile routes
		api.GET("/profile", userHandler.GetProfile)
		api.PUT("/profile", userHandler.UpdateProfile)
		api.DELETE("/profile", userHandler.DeleteProfile)

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.RequireRole("admin", log))
		{
			admin.GET("/users", userHandler.ListUsers)
			admin.POST("/users", userHandler.CreateUser)
			admin.PUT("/users/:id", userHandler.UpdateUser)
			admin.DELETE("/users/:id", userHandler.DeleteUser)
		}
	}
}

// Gin Auth API
// @title Gin Authentication API
// @version 1.0
// @description A production-ready authentication API built with Gin framework, featuring JWT authentication, RBAC, and comprehensive user management
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/shahariaz/gin-auth-service/internal/config"
	"github.com/shahariaz/gin-auth-service/internal/database"
	"github.com/shahariaz/gin-auth-service/internal/logger"
	"github.com/shahariaz/gin-auth-service/internal/middleware"
	"github.com/shahariaz/gin-auth-service/internal/router"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	mredis "github.com/ulule/limiter/v3/drivers/store/redis"
	
	_ "github.com/shahariaz/gin-auth-service/docs" // Import generated docs
)

func main() {
	// Parse flags
	env := flag.String("env", "development", "Environment: development/production")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found or error loading: %v. Using system env vars.", err)

	}
	cfg := config.LoadConfig(*env)

	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Logger setup
	log := logger.NewLogger(cfg.GinMode)

	// Database setup
	db, err := database.NewDatabase(cfg.DB_DSN)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	if err := database.RunMigrations(db, log); err != nil {
		log.Fatal("Failed to run migrations: ", err)
	}

	// Gin engine
	r := gin.New()

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rate limiting with ulule/limiter (Redis for prod, in-memory for dev)
	var rateStore limiter.Store
	if cfg.RedisURL != "" {
		redisClient := redis.NewClient(&redis.Options{Addr: cfg.RedisURL})
		rateStore, err = mredis.NewStore(redisClient)
		if err != nil {
			log.Fatal("Failed to initialize Redis store: ", err)
		}
	} else {
		rateStore = memory.NewStore() // In-memory for dev
	}
	rate, err := limiter.NewRateFromFormatted("10-S") // 10 requests per second
	if err != nil {
		log.Fatal("Failed to parse rate limit: ", err)
	}
	rateLimiter := mgin.NewMiddleware(limiter.New(rateStore, rate, limiter.WithTrustForwardHeader(true)))
	r.Use(rateLimiter)
	r.Use(func(c *gin.Context) {
		if c.Writer.Status() == http.StatusTooManyRequests {
			log.WithFields(map[string]interface{}{
				"ip": c.ClientIP(),
			}).Warn("Rate limit exceeded")
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests",
				"retry_after": c.GetHeader("X-RateLimit-Reset"),
			})
			c.Abort()
			return
		}
		c.Next()
	})

	// Other middleware
	r.Use(middleware.LoggingMiddleware(log))
	r.Use(middleware.TimeoutMiddleware(5*time.Second, log))
	r.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	})
	r.Use(middleware.ErrorMiddleware(log))

	// Static files
	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "./static/favicon.ico")

	// Health check
	// @Summary Health check
	// @Description Check if the API is running and healthy
	// @Tags System
	// @Produce json
	// @Success 200 {object} map[string]interface{} "API is healthy"
	// @Router /health [get]
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "version": cfg.AppVersion})
	})

	// Setup routes
	router.SetupRoutes(r, cfg, db, log)

	// Start server (HTTP for local development - no TLS)
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		// Use ListenAndServe for HTTP
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed: ", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced shutdown: ", err)
	}
	log.Info("Server exited gracefully")
}

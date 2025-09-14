package config

import (
	"log" // For warnings (or use your logger)
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppVersion      string
	GinMode         string
	Port            string
	JWT_SECRET      []byte
	AllowOrigins    []string
	RateLimitPerSec int
	DB_DSN          string
	RedisURL        string
}

func LoadConfig(env string) *Config { // Changed to return pointer for consistency
	cfg := &Config{ // Use pointer
		AppVersion:      "1.0.0",
		GinMode:         strings.TrimSpace(os.Getenv("GIN_MODE")), // Read from env
		Port:            strings.TrimSpace(os.Getenv("PORT")), // Trim whitespace
		JWT_SECRET:      []byte(os.Getenv("JWT_SECRET")),
		AllowOrigins:    []string{"http://localhost:3000"},
		RateLimitPerSec: 10,
		DB_DSN:          strings.TrimSpace(os.Getenv("DB_DSN")),    // Trim
		RedisURL:        strings.TrimSpace(os.Getenv("REDIS_URL")), // Trim
	}

	// Set default GIN_MODE if not provided
	if cfg.GinMode == "" {
		if env == "development" {
			cfg.GinMode = "debug"
		} else {
			cfg.GinMode = "release"
		}
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	if env == "production" {
		cfg.GinMode = "release"
		if cfg.Port == "8080" { // Default if unset
			cfg.Port = os.Getenv("PORT")
			if cfg.Port == "" {
				cfg.Port = "8080"
			}
		}
		if origins := os.Getenv("ALLOWED_ORIGINS"); origins != "" {
			trimmedOrigins := strings.Split(strings.TrimSpace(origins), ",")
			for i := range trimmedOrigins {
				trimmedOrigins[i] = strings.TrimSpace(trimmedOrigins[i]) // Trim each origin
			}
			cfg.AllowOrigins = trimmedOrigins
		}
		rateStr := os.Getenv("RATE_LIMIT_PER_SEC")
		if rateStr != "" {
			if rate, err := strconv.Atoi(rateStr); err == nil && rate > 0 {
				cfg.RateLimitPerSec = rate
			} else if err != nil {
				log.Printf("Warning: Invalid RATE_LIMIT_PER_SEC '%s': %v. Using default %d", rateStr, err, cfg.RateLimitPerSec)
			}
		}
		if cfg.RedisURL == "" {
			log.Println("Warning: REDIS_URL not set in production. Token blacklisting disabled.")
		}
	}

	// Validate required fields
	if string(cfg.JWT_SECRET) == "" {
		panic("JWT_SECRET not set - required for JWT signing")
	}
	if cfg.DB_DSN == "" {
		panic("DB_DSN not set - required for database connection")
	}
	if cfg.Port == "" {
		panic("PORT not set - required for server binding")
	}

	return cfg
}

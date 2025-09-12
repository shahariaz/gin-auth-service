package database

import (
	"github.com/shahariaz/gin-auth-service/internal/model"
	"github.com/sirupsen/logrus"
)

func RunMigrations(db *Database, log *logrus.Logger) error {
	log.Info("Running database migrations...")
	
	// Run auto migrations
	if err := db.AutoMigrate(&model.Role{}, &model.User{}); err != nil {
		log.WithError(err).Error("Failed to run auto migrations")
		return err
	}
	
	// Create default roles if they don't exist
	var userRole model.Role
	if err := db.Where("name = ?", "user").First(&userRole).Error; err != nil {
		userRole = model.Role{
			Name:        "user",
			Description: "Regular user role",
		}
		if err := db.Create(&userRole).Error; err != nil {
			log.WithError(err).Error("Failed to create user role")
			return err
		}
		log.Info("Created default user role")
	}
	
	var adminRole model.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		adminRole = model.Role{
			Name:        "admin", 
			Description: "Administrator role",
		}
		if err := db.Create(&adminRole).Error; err != nil {
			log.WithError(err).Error("Failed to create admin role")
			return err
		}
		log.Info("Created default admin role")
	}
	
	log.Info("Database migrations completed successfully")
	return nil
}

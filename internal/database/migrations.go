package database

import (
	"github.com/shahariaz/gin-auth-service/internal/model"
	"github.com/sirupsen/logrus"
)

func RunMigrations(db *Database, log *logrus.Logger) error {
	log.Info("Running database migrations...")
	return db.AutoMigrate(&model.User{}, &model.Role{})
}

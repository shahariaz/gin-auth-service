package service

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/shahariaz/gin-auth-service/internal/database"
	"github.com/shahariaz/gin-auth-service/internal/model"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	db        *database.Database
	validator *validator.Validate
	log       *logrus.Logger
}

func NewUserService(db *database.Database, validator *validator.Validate, log *logrus.Logger) *UserService {
	return &UserService{db: db, validator: validator, log: log}
}

func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := s.db.Preload("Role").Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) UpdateUserProfile(username, email string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	user.Email = email
	user.UpdatedAt = time.Now()
	if err := s.validator.Struct(user); err != nil {
		return nil, err
	}
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) DeleteUser(username string) error {
	return s.db.Where("username = ?", username).Delete(&model.User{}).Error
}

func (s *UserService) ListUsers() ([]model.User, error) {
	var users []model.User
	if err := s.db.Preload("Role").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) CreateUser(user *model.User) error {
	if err := s.validator.Struct(user); err != nil {
		return err
	}
	var existing model.User
	if err := s.db.Where("email = ? OR username = ?", user.Email, user.Username).First(&existing).Error; err == nil {
		return errors.New("user already exists")
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return s.db.Create(user).Error
}

func (s *UserService) UpdateUser(id uint, username, email string, roleID uint) (*model.User, error) {
	var user model.User
	if err := s.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	user.Username = username
	user.Email = email
	user.RoleID = roleID
	user.UpdatedAt = time.Now()
	if err := s.validator.Struct(user); err != nil {
		return nil, err
	}
	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) DeleteUserByID(id uint) error {
	return s.db.Where("id = ?", id).Delete(&model.User{}).Error
}

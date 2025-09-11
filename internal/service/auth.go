package service

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/shahariaz/gin-auth-service/internal/database"
	"github.com/shahariaz/gin-auth-service/internal/lib"
	"github.com/shahariaz/gin-auth-service/internal/model"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db         *database.Database
	validator  *validator.Validate
	tokenStore lib.TokenStore
	secret     []byte
	log        *logrus.Logger
}

func NewAuthService(db *database.Database, validator *validator.Validate, tokenStore lib.TokenStore, secret []byte, log *logrus.Logger) *AuthService {
	return &AuthService{db: db, validator: validator, tokenStore: tokenStore, secret: secret, log: log}
}

func (s *AuthService) Register(user *model.User, password string) error {
	if err := s.validator.Struct(user); err != nil {
		return err
	}

	// Check if user exists
	var existing model.User
	if err := s.db.Where("email = ? OR username = ?", user.Email, user.Username).First(&existing).Error; err == nil {
		return errors.New("user already exists")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.db.Create(user).Error
}

func (s *AuthService) Login(email, password string) (*model.User, string, string, error) {
	var user model.User
	if err := s.db.Preload("Role").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, "", "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", "", errors.New("invalid password")
	}

	accessToken, err := lib.GenerateAccessToken(user.ID, user.Username, user.Role.Name, s.secret)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := lib.GenerateRefreshToken(user.ID, user.Username, s.secret)
	if err != nil {
		return nil, "", "", err
	}

	return &user, accessToken, refreshToken, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	isBlacklisted, err := s.tokenStore.IsBlacklisted(refreshToken)
	if err != nil {
		return "", err
	}
	if isBlacklisted {
		return "", errors.New("refresh token blacklisted")
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["username"] == nil || claims["user_id"] == nil {
		return "", errors.New("invalid token claims")
	}

	var user model.User
	if err := s.db.Preload("Role").Where("id = ?", uint(claims["user_id"].(float64))).First(&user).Error; err != nil {
		return "", errors.New("user not found")
	}

	accessToken, err := lib.GenerateAccessToken(user.ID, user.Username, user.Role.Name, s.secret)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *AuthService) Logout(refreshToken string) error {
	return s.tokenStore.Blacklist(refreshToken, 7*24*time.Hour)
}

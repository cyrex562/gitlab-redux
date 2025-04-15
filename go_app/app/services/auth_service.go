package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/gitlab-org/gitlab-redux/app/models"
	"gorm.io/gorm"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// AuthService handles authentication operations
type AuthService struct {
	db         *gorm.DB
	jwtSecret  []byte
	tokenTTL   time.Duration
}

// NewAuthService creates a new AuthService
func NewAuthService(db *gorm.DB, jwtSecret string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		db:         db,
		jwtSecret:  []byte(jwtSecret),
		tokenTTL:   tokenTTL,
	}
}

// ValidateToken validates a JWT token and returns the associated user
func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Validate claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, ErrExpiredToken
			}
		}

		// Get user ID from claims
		if userID, ok := claims["user_id"].(float64); ok {
			var user models.User
			if err := s.db.First(&user, uint(userID)).Error; err != nil {
				return nil, err
			}
			return &user, nil
		}
	}

	return nil, ErrInvalidToken
}

// GenerateToken generates a JWT token for a user
func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	// Create claims
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(s.tokenTTL).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	return token.SignedString(s.jwtSecret)
} 
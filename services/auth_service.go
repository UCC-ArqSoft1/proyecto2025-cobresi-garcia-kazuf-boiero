package services

import (
    "errors"
    "fmt"

    "github.com/alesio/gestion-actividades-deportivas/config"
    "github.com/alesio/gestion-actividades-deportivas/models"
    "github.com/golang-jwt/jwt/v5"
    "gorm.io/gorm"
)

// JWTClaims represents the claims encoded inside application tokens.
type JWTClaims struct {
    UserID uint   `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

// AuthService coordinates authentication and token management logic.
type AuthService struct {
    db  *gorm.DB
    cfg *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
    return &AuthService{db: db, cfg: cfg}
}

// Authenticate validates the provided credentials and returns the matching user when valid.
func (s *AuthService) Authenticate(email, password string) (*models.User, error) {
    var user models.User
    if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("invalid credentials")
        }
        return nil, err
    }

    // TODO: Compare password hash with provided password and return error on mismatch.

    return &user, nil
}

// GenerateJWT should issue a signed JWT for the given user.
func (s *AuthService) GenerateJWT(user *models.User) (string, error) {
    // TODO: Build JWT claims, sign with cfg.JWTSecret, and return serialized token.
    return "", fmt.Errorf("token generation not implemented")
}

// ValidateJWT should parse and validate the token string.
func (s *AuthService) ValidateJWT(token string) (*JWTClaims, error) {
    // TODO: Parse token with jwt.ParseWithClaims using cfg.JWTSecret and return claims on success.
    return nil, fmt.Errorf("token validation not implemented")
}

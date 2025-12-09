package service

import (
	"context"
	"errors"
	"server/internal/core/domain"
	"server/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	userRepo  repository.Repository[domain.User]
	jwtSecret []byte
}

func NewAuthService(userRepo repository.Repository[domain.User], secret string) AuthService {
	return &AuthServiceImpl{
		userRepo:  userRepo,
		jwtSecret: []byte(secret),
	}
}

func (s *AuthServiceImpl) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.FindOne(ctx, "email = ?", email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if user.PasswordHash == nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Access Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	accessToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Refresh Token (Simplified for now)
	refreshToken := "refresh_token_placeholder"

	return accessToken, refreshToken, nil
}

func (s *AuthServiceImpl) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	return "", "", nil
}

func (s *AuthServiceImpl) Logout(ctx context.Context, tokenString string) error {
	return nil
}

func (s *AuthServiceImpl) RegisterStaff(ctx context.Context, user *domain.User, plainPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	hash := string(hashedPassword)
	user.PasswordHash = &hash
	user.IsActive = true

	return s.userRepo.Create(ctx, user)
}

func (s *AuthServiceImpl) RequestPasswordReset(ctx context.Context, email string) error {
	return nil
}

func (s *AuthServiceImpl) ConfirmPasswordReset(ctx context.Context, token, newPassword string) error {
	return nil
}

func (s *AuthServiceImpl) ResetPassword(ctx context.Context, userID int, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	hash := string(hashedPassword)

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.PasswordHash = &hash
	return s.userRepo.Update(ctx, user)
}

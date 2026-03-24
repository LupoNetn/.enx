package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/luponetn/enx/internal/config"
	"github.com/luponetn/enx/internal/db"
	"github.com/luponetn/enx/internal/utils"
)

type AuthService struct {
	queries *db.Queries
	config  *config.Config
}

func NewAuthService(queries *db.Queries, cfg *config.Config) *AuthService {
	return &AuthService{queries: queries, config: cfg}
}


func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	// check if email already exists
	_, err := s.queries.GetUserByEmail(ctx, input.Email)
	if err == nil {
		return nil, errors.New("email already in use")
	}

	// hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return nil, errors.New("internal server error")
	}

	// create user
	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Email:    input.Email,
		Name:     input.Name,
		Password: hashedPassword,
	})
	if err != nil {
		slog.Error("failed to create user", "error", err)
		return nil, errors.New("failed to create user")
	}

	// generate tokens
	tokens, err := GenerateTokenPair(user.ID, s.config.JWTAccessSecret, s.config.JWTRefreshSecret)
	if err != nil {
		slog.Error("failed to generate tokens", "error", err)
		return nil, errors.New("failed to generate tokens")
	}

	return &AuthResponse{User: user, Tokens: tokens}, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	// get user with password
	user, err := s.queries.GetUserForAuth(ctx, input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// compare password
	if !utils.ComparePassword(user.Password, input.Password) {
		return nil, errors.New("invalid email or password")
	}

	// generate tokens
	tokens, err := GenerateTokenPair(user.ID, s.config.JWTAccessSecret, s.config.JWTRefreshSecret)
	if err != nil {
		slog.Error("failed to generate tokens", "error", err)
		return nil, errors.New("failed to generate tokens")
	}

	return &AuthResponse{
		User: map[string]any{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
		Tokens: tokens,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// validate using REFRESH secret
	claims, err := ValidateToken(refreshToken, s.config.JWTRefreshSecret)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	userID, err := StringToUUID(claims.UserID)
	if err != nil {
		return nil, errors.New("invalid token claims")
	}

	// verify user still exists
	_, err = s.queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user no longer exists")
	}

	// issue fresh token pair
	tokens, err := GenerateTokenPair(userID, s.config.JWTAccessSecret, s.config.JWTRefreshSecret)
	if err != nil {
		slog.Error("failed to generate tokens", "error", err)
		return nil, errors.New("failed to generate tokens")
	}

	return tokens, nil
}

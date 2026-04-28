package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/AdelmoMJunior/GoACBr/internal/auth"
	"github.com/AdelmoMJunior/GoACBr/internal/domain"
	"github.com/AdelmoMJunior/GoACBr/internal/dto"
	"github.com/AdelmoMJunior/GoACBr/internal/repository"
	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
)

type AuthService interface {
	Login(ctx context.Context, req *dto.LoginRequest, ipAddress, userAgent string) (*dto.LoginResponse, error)
	Refresh(ctx context.Context, req *dto.RefreshRequest, ipAddress, userAgent string) (*dto.LoginResponse, error)
	Logout(ctx context.Context, sessionID uuid.UUID, jti string, expiresAt time.Time) error
}

type authService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	tokenSvc    *auth.TokenService
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, tokenSvc *auth.TokenService) AuthService {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		tokenSvc:    tokenSvc,
	}
}

func (s *authService) Login(ctx context.Context, req *dto.LoginRequest, ipAddress, userAgent string) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, apperror.NewUnauthorized("invalid email or password")
	}

	if !user.IsActive {
		return nil, apperror.NewUnauthorized("account is disabled")
	}

	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, apperror.NewUnauthorized("invalid email or password")
	}

	sessionID := uuid.New()
	tokens, err := s.tokenSvc.GenerateTokenPair(user.ID, user.Email, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	refreshHash, err := auth.GenerateRandomHash(tokens.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to hash refresh token: %w", err)
	}

	session := &domain.Session{
		ID:               sessionID,
		UserID:           user.ID,
		RefreshTokenHash: refreshHash,
		IPAddress:        ipAddress,
		UserAgent:        userAgent,
		IsRevoked:        false,
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour), // Must match config refresh TTL
	}

	if err := s.sessionRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &dto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600, // Hardcoded for response, must match config access TTL
	}, nil
}

func (s *authService) Refresh(ctx context.Context, req *dto.RefreshRequest, ipAddress, userAgent string) (*dto.LoginResponse, error) {
	claims, err := s.tokenSvc.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, apperror.NewUnauthorized("invalid refresh token")
	}

	sessionID, err := uuid.Parse(claims.ID)
	if err != nil {
		return nil, apperror.NewUnauthorized("invalid session id in token")
	}

	session, err := s.sessionRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, apperror.NewUnauthorized("session not found")
	}

	if session.IsRevoked || session.ExpiresAt.Before(time.Now()) {
		return nil, apperror.NewUnauthorized("session expired or revoked")
	}

	if !auth.CheckPasswordHash(req.RefreshToken, session.RefreshTokenHash) {
		return nil, apperror.NewUnauthorized("invalid refresh token")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, apperror.NewUnauthorized("invalid user id in token")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || !user.IsActive {
		return nil, apperror.NewUnauthorized("user not found or inactive")
	}

	// Revoke old session and create a new one (Rotation)
	_ = s.sessionRepo.RevokeSession(ctx, sessionID)

	newSessionID := uuid.New()
	tokens, err := s.tokenSvc.GenerateTokenPair(user.ID, user.Email, newSessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new tokens: %w", err)
	}

	refreshHash, err := auth.GenerateRandomHash(tokens.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to hash new refresh token: %w", err)
	}

	newSession := &domain.Session{
		ID:               newSessionID,
		UserID:           user.ID,
		RefreshTokenHash: refreshHash,
		IPAddress:        ipAddress,
		UserAgent:        userAgent,
		IsRevoked:        false,
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.sessionRepo.CreateSession(ctx, newSession); err != nil {
		return nil, fmt.Errorf("failed to create new session: %w", err)
	}

	return &dto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	}, nil
}

func (s *authService) Logout(ctx context.Context, sessionID uuid.UUID, jti string, expiresAt time.Time) error {
	_ = s.sessionRepo.RevokeSession(ctx, sessionID)
	return s.sessionRepo.BlacklistToken(ctx, jti, expiresAt)
}

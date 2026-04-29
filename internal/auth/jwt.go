package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// Claims represents the standard JWT claims plus custom fields.
type Claims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	CompanyID string `json:"company_id,omitempty"` // For operations scoped to a specific company
	jwt.RegisteredClaims
}

// TokenService handles JWT creation and validation.
type TokenService struct {
	secretKey       []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewTokenService creates a new TokenService.
func NewTokenService(secret string, accessTTL, refreshTTL time.Duration) *TokenService {
	return &TokenService{
		secretKey:       []byte(secret),
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

// GenerateTokenPair creates both an access token and a refresh token.
// The sessionID should be stored in the refresh token so we can tie it to a DB session.
func (s *TokenService) GenerateTokenPair(userID uuid.UUID, email string, sessionID uuid.UUID) (*TokenPair, error) {
	now := time.Now()

	// 1. Access Token
	accessClaims := Claims{
		UserID: userID.String(),
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   userID.String(),
			ID:        uuid.New().String(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString(s.secretKey)
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	// 2. Refresh Token
	// The Subject is the UserID, the ID is the SessionID
	refreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Subject:   userID.String(),
		ID:        sessionID.String(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString(s.secretKey)
	if err != nil {
		return nil, fmt.Errorf("sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
	}, nil
}

// ValidateToken parses and validates a token string.
func (s *TokenService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// ValidateRefreshToken parses and validates a refresh token string.
func (s *TokenService) ValidateRefreshToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token claims")
	}

	return claims, nil
}

// HashPassword generates a bcrypt hash of the given password.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// CheckPasswordHash compares a password with a hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomHash generates a SHA256 hash of the data (useful for long tokens).
func GenerateRandomHash(data string) (string, error) {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:]), nil
}

// CheckTokenHash compares a token with its SHA256 hash.
func CheckTokenHash(token, hash string) bool {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:]) == hash
}

// --- Context keys and helpers ---

type contextKey string

const (
	userIDKey    contextKey = "user_id"
	companyIDKey contextKey = "company_id"
	jtiKey       contextKey = "jti"
	sessionIDKey contextKey = "session_id"
	expiresAtKey contextKey = "expires_at"
)

// WithUserID adds a user ID to the context.
func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID retrieves the user ID from the context.
func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	val, ok := ctx.Value(userIDKey).(uuid.UUID)
	return val, ok
}

// WithCompanyID adds a company ID to the context.
func WithCompanyID(ctx context.Context, companyID uuid.UUID) context.Context {
	return context.WithValue(ctx, companyIDKey, companyID)
}

// GetCompanyID retrieves the company ID from the context.
func GetCompanyID(ctx context.Context) (uuid.UUID, bool) {
	val, ok := ctx.Value(companyIDKey).(uuid.UUID)
	return val, ok
}

// WithClaims adds the JWT claims (jti, sessionID, expiresAt) to the context.
func WithClaims(ctx context.Context, jti string, sessionID uuid.UUID, expiresAt time.Time) context.Context {
	ctx = context.WithValue(ctx, jtiKey, jti)
	ctx = context.WithValue(ctx, sessionIDKey, sessionID)
	ctx = context.WithValue(ctx, expiresAtKey, expiresAt)
	return ctx
}

// GetJTI retrieves the JWT ID (jti) from the context.
func GetJTI(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(jtiKey).(string)
	return val, ok
}

// GetSessionID retrieves the session ID from the context.
func GetSessionID(ctx context.Context) (uuid.UUID, bool) {
	val, ok := ctx.Value(sessionIDKey).(uuid.UUID)
	return val, ok
}

// GetExpiresAt retrieves the token expiration time from the context.
func GetExpiresAt(ctx context.Context) (time.Time, bool) {
	val, ok := ctx.Value(expiresAtKey).(time.Time)
	return val, ok
}

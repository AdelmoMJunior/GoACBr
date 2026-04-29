package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/AdelmoMJunior/GoACBr/internal/auth"
	"github.com/AdelmoMJunior/GoACBr/internal/repository"
	"github.com/AdelmoMJunior/GoACBr/pkg/apperror"
	"github.com/AdelmoMJunior/GoACBr/pkg/httputil"
)

// AuthMiddleware validates the JWT token and adds UserID to context.
func AuthMiddleware(tokenSvc *auth.TokenService, sessionRepo repository.SessionRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				httputil.SendError(w, apperror.NewUnauthorized("missing or invalid authorization header"))
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := tokenSvc.ValidateToken(tokenStr)
			if err != nil {
				httputil.SendError(w, apperror.NewUnauthorized("invalid token"))
				return
			}

			// Check blacklist
			blacklisted, err := sessionRepo.IsTokenBlacklisted(r.Context(), claims.ID)
			if err != nil || blacklisted {
				httputil.SendError(w, apperror.NewUnauthorized("token has been revoked"))
				return
			}

			userID, err := uuid.Parse(claims.UserID)
			if err != nil {
				httputil.SendError(w, apperror.NewUnauthorized("invalid user id in token"))
				return
			}

			// Inject user ID into context
			ctx := auth.WithUserID(r.Context(), userID)

			// Inject JTI, SessionID and ExpiresAt for Logout support
			sessionID := uuid.Nil
			if claims.ID != "" {
				if sid, parseErr := uuid.Parse(claims.ID); parseErr == nil {
					sessionID = sid
				}
			}
			var expiresAt time.Time
			if claims.ExpiresAt != nil {
				expiresAt = claims.ExpiresAt.Time
			}
			ctx = auth.WithClaims(ctx, claims.ID, sessionID, expiresAt)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CompanyGuard ensures the user is linked to the requested company.
// Requires the route to have {companyID} param. (We'll use go-chi to extract it)
func CompanyGuard(compRepo repository.CompanyRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// In go-chi, we get URL params via chi.URLParam
			// To keep middleware generic, we can extract from a known header or chi context
			// We'll rely on a custom context extractor or chi later.
			// For now, assume it's extracted and placed in context or passed via header X-Company-ID
			compIDStr := r.Header.Get("X-Company-ID")
			if compIDStr == "" {
				httputil.SendError(w, apperror.NewBadRequest("missing X-Company-ID header"))
				return
			}

			companyID, err := uuid.Parse(compIDStr)
			if err != nil {
				httputil.SendError(w, apperror.NewBadRequest("invalid company id format"))
				return
			}

			userID, ok := auth.GetUserID(r.Context())
			if !ok {
				httputil.SendError(w, apperror.NewUnauthorized("user not authenticated"))
				return
			}

			// Check linkage
			users, err := compRepo.GetUsersByCompany(r.Context(), companyID)
			if err != nil {
				httputil.SendError(w, apperror.NewInternal(errors.New("failed to verify access")))
				return
			}

			hasAccess := false
			for _, u := range users {
				if u.ID == userID {
					hasAccess = true
					break
				}
			}

			if !hasAccess {
				httputil.SendError(w, apperror.NewForbidden("you do not have access to this company"))
				return
			}

			ctx := auth.WithCompanyID(r.Context(), companyID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CertValidator ensures the company has a valid, non-expired certificate.
func CertValidator(certRepo repository.CertificateRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			companyID, ok := auth.GetCompanyID(r.Context())
			if !ok {
				httputil.SendError(w, apperror.NewInternal(errors.New("company id missing in context")))
				return
			}

			cert, err := certRepo.GetByCompanyID(r.Context(), companyID)
			if err != nil {
				httputil.SendError(w, apperror.NewBadRequest("company must have a digital certificate to perform this operation"))
				return
			}

			// In a real app, verify expiration
			if cert.ValidUntil.Before(time.Now()) {
				httputil.SendError(w, apperror.NewBadRequest("digital certificate has expired"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

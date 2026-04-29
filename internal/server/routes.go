package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/AdelmoMJunior/GoACBr/internal/auth"
	"github.com/AdelmoMJunior/GoACBr/internal/handler"
	appmiddleware "github.com/AdelmoMJunior/GoACBr/internal/middleware"
	"github.com/AdelmoMJunior/GoACBr/internal/repository"

	_ "github.com/AdelmoMJunior/GoACBr/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// SetupRoutes configures all HTTP routes.
func SetupRoutes(
	r chi.Router,
	tokenSvc *auth.TokenService,
	sessionRepo repository.SessionRepository,
	authH *handler.AuthHandler,
	userH *handler.UserHandler,
	compH *handler.CompanyHandler,
	certH *handler.CertificateHandler,
	nfeH *handler.NFeHandler,
	distH *handler.DistributionHandler,
	invH *handler.InvoiceHandler,
	healthH *handler.HealthHandler,
) {
	// Standard middlewares
	r.Use(middleware.RealIP)
	r.Use(appmiddleware.RequestID)
	r.Use(appmiddleware.Logger)
	r.Use(appmiddleware.Recovery)

	// CORS config
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Company-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// API v1 router
	r.Route("/api/v1", func(api chi.Router) {
		// Public routes
		healthH.RegisterRoutes(api)
		authH.RegisterRoutes(api) // Login/Refresh
		userH.RegisterRoutes(api) // Register

		// Swagger UI
		api.Get("/swagger/*", httpSwagger.WrapHandler)

		// Protected routes
		api.Group(func(protected chi.Router) {
			protected.Use(appmiddleware.AuthMiddleware(tokenSvc, sessionRepo))
			
			// Routes that don't need company context
			// (authH logout, userH me routes)
			authH.RegisterProtectedRoutes(protected)
			userH.RegisterProtectedRoutes(protected)
			compH.RegisterRoutes(protected)

			// The following need CompanyGuard, which is applied inside their own RegisterRoutes
			certH.RegisterRoutes(protected)
			nfeH.RegisterRoutes(protected)
			distH.RegisterRoutes(protected)
			invH.RegisterRoutes(protected)
		})
	})
}

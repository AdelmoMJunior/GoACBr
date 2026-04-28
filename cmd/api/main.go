package main

import (
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"

	"github.com/AdelmoMJunior/GoACBr/internal/acbr"
	"github.com/AdelmoMJunior/GoACBr/internal/auth"
	"github.com/AdelmoMJunior/GoACBr/internal/config"
	"github.com/AdelmoMJunior/GoACBr/internal/crypto"
	"github.com/AdelmoMJunior/GoACBr/internal/handler"
	"github.com/AdelmoMJunior/GoACBr/internal/repository"
	"github.com/AdelmoMJunior/GoACBr/internal/server"
	"github.com/AdelmoMJunior/GoACBr/internal/service"
	"github.com/AdelmoMJunior/GoACBr/internal/storage"
	"github.com/AdelmoMJunior/GoACBr/pkg/logger"
)

func main() {
	// 1. Config & Logger
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}
	logger.Setup(cfg.Log.Level, cfg.Log.Format)

	// 2. Database
	db, err := repository.NewDB(cfg.Database)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// 3. Redis / Cache
	redisCache, err := repository.NewRedisCache(cfg.Redis)
	if err != nil {
		slog.Warn("Cache initialization issue", "error", err)
	} else {
		defer redisCache.Close()
	}

	// 4. Storage (B2)
	storageProv, err := storage.NewB2Storage(cfg.B2)
	if err != nil {
		slog.Error("Failed to initialize B2 storage", "error", err)
		// We could exit or run in degraded mode, let's exit for scaffold purity
		os.Exit(1)
	}

	// 5. ACBrLib Pool
	pool, err := acbr.NewHandlePool(
		5, // PoolSize
		cfg.ACBr.LibPath,
		"/app/acbr.ini", // ConfigPath
		"",              // CryptKey
	)
	if err != nil {
		slog.Error("Failed to initialize ACBrLib pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// 6. Security/Crypto
	cryptoSvc, err := crypto.NewAESService(cfg.Encryption.Key)
	if err != nil {
		slog.Error("Failed to init crypto", "error", err)
		os.Exit(1)
	}

	tokenSvc := auth.NewTokenService(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)

	// 7. Repositories
	userRepo := repository.NewUserRepository(db)
	compRepo := repository.NewCompanyRepository(db)
	certRepo := repository.NewCertificateRepository(db)
	sessRepo := repository.NewSessionRepository(db)
	invRepo := repository.NewInvoiceRepository(db)
	distRepo := repository.NewDistributionRepository(db)
	// auditRepo := repository.NewAuditRepository(db)

	// 8. Services
	authSvc := service.NewAuthService(userRepo, sessRepo, tokenSvc)
	userSvc := service.NewUserService(userRepo)
	compSvc := service.NewCompanyService(compRepo, certRepo, cryptoSvc)
	certSvc := service.NewCertificateService(certRepo, cryptoSvc)
	nfeSvc := service.NewNFeService(compRepo, certRepo, invRepo, pool, storageProv, cryptoSvc)
	evtSvc := service.NewEventService(compRepo, certRepo, invRepo, pool, cryptoSvc)
	distSvc := service.NewDistributionService(compRepo, certRepo, distRepo, pool, cryptoSvc)
	invSvc := service.NewInvoiceService(invRepo)
	// emailSvc := service.NewEmailService(compRepo, cryptoSvc)

	// 9. Handlers
	authH := handler.NewAuthHandler(authSvc)
	userH := handler.NewUserHandler(userSvc)
	compH := handler.NewCompanyHandler(compSvc, compRepo)
	certH := handler.NewCertificateHandler(certSvc, compRepo)
	nfeH := handler.NewNFeHandler(nfeSvc, evtSvc, compRepo, certRepo)
	distH := handler.NewDistributionHandler(distSvc, compRepo, certRepo)
	invH := handler.NewInvoiceHandler(invSvc, compRepo)
	healthH := handler.NewHealthHandler()

	// 10. Router
	r := chi.NewRouter()
	server.SetupRoutes(r, tokenSvc, sessRepo, authH, userH, compH, certH, nfeH, distH, invH, healthH)

	// 11. Start Server
	srv := server.NewServer(cfg.Server, r)
	srv.Start()
}

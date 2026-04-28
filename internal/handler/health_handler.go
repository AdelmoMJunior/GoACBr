package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AdelmoMJunior/GoACBr/pkg/httputil"
)

type HealthHandler struct {
	// Can inject DB/Redis pingers here to check readiness
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) RegisterRoutes(r chi.Router) {
	r.Get("/health", h.Liveness)
	r.Get("/ready", h.Readiness)
}

func (h *HealthHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	httputil.SendJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	// In a real app, check DB, Redis, ACBrLib pool
	httputil.SendJSON(w, http.StatusOK, map[string]string{
		"status": "ready",
		"db":     "ok",
		"redis":  "ok",
		"acbr":   "ok",
	})
}

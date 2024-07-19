package handler

import (
	"net/http"

	"github.com/ecol-master/sharing-wh-machines/internal/config"
	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/ecol-master/sharing-wh-machines/internal/http/middlewares"
	"github.com/ecol-master/sharing-wh-machines/internal/service"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	service *service.Service
	cfg     *config.Config
}

func New(db *sqlx.DB, cfg *config.Config) *Handler {
	return &Handler{
		service: service.New(db),
		cfg:     cfg,
	}
}

func (h *Handler) MakeHTTPHandler() http.Handler {
	mux := http.NewServeMux()

	// api methods for web-application
	mux.Handle("GET /get_all_users", h.makeAdminHandler(h.GetAllUsers))
	mux.Handle("GET /get_user", h.makeAdminHandler(h.GetUserByID))

	mux.Handle("GET /get_all_machines", h.makeAdminHandler(h.GetAllMachines))
	mux.Handle("GET /get_machine", h.makeAdminHandler(h.GetMachineByID))

	mux.Handle("GET /get_all_sessions", h.makeAdminHandler(h.GetAllSessions))
	mux.Handle("GET /get_session", h.makeAdminHandler(h.GetSessionByID))

	// auth
	mux.HandleFunc("POST /login", h.Login)

	// Lock, Unlock, Pause handler
	mux.Handle("POST /unlock_machine", h.makeWorkerHandler(h.UnlockMachine))
	mux.Handle("POST /lock_machine", h.makeWorkerHandler(h.LockMachine))

	// handler to register (or make active after failed) arduino in system
	mux.Handle("POST /register_machine", http.HandlerFunc(h.RegisterMachine))

	// logging all request with LoggingMiddleware
	return middlewares.LoggingMiddleware(mux)
}

// function making handler from RoleBasedAccess middleware with entities.Admin role
func (h *Handler) makeAdminHandler(handleFunc func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return middlewares.RoleBasedAccess(h.cfg.Secret, entities.Admin, http.HandlerFunc(handleFunc))
}

func (h *Handler) makeWorkerHandler(handleFunc func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return middlewares.RoleBasedAccess(h.cfg.Secret, entities.Worker, http.HandlerFunc(handleFunc))
}

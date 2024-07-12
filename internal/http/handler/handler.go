package handler

import (
	"net/http"

	"github.com/ecol-master/sharing-wh-machines/internal/config"
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

	// api methods
	mux.Handle("GET /get_all_users", h.makeAdminHandler(h.GetAllUsers))
	mux.Handle("GET /get_user", h.makeAdminHandler(h.GetUserByID))

	mux.Handle("GET /get_all_machines", h.makeAdminHandler(h.GetAllMachines))
	mux.Handle("GET /get_machine", h.makeAdminHandler(h.GetMachineByID))

	mux.Handle("GET /get_all_sessions", h.makeAdminHandler(h.GetAllSessions))
	mux.Handle("GET /get_session", h.makeAdminHandler(h.GetSessionByID))

	// auth
	mux.HandleFunc("POST /login", h.Login)

	// logging all request with LoggingMiddleware
	return middlewares.LoggingMiddleware(mux)
}

// function making handler with IsAdmin middleware
func (h *Handler) makeAdminHandler(handleFunc func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return middlewares.IsAdmin(h.cfg.Secret, http.HandlerFunc(handleFunc))
}

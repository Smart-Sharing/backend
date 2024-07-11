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
	mux.HandleFunc("GET /get_all_users", h.GetAllUsers)
	mux.HandleFunc("GET /get_user", h.GetUserByID)

	mux.HandleFunc("GET /get_all_machines", h.GetAllMachines)
	mux.HandleFunc("GET /get_machine", h.GetMachineByID)

	mux.HandleFunc("GET /get_all_sessions", h.GetAllSessions)
	mux.HandleFunc("GET /get_session", h.GetSessionByID)

	// auth
	mux.HandleFunc("POST /login", h.Login)

	// logging all request with LoggingMiddleware
	return middlewares.LoggingMiddleware(mux)
}

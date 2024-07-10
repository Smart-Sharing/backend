package handler

import (
	"net/http"

	"github.com/ecol-master/sharing-wh-machines/internal/http/middlewares"
	"github.com/ecol-master/sharing-wh-machines/internal/service"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	service *service.Service
}

func New(db *sqlx.DB) *Handler {
	return &Handler{
		service: service.New(db),
	}
}

func (h *Handler) MakeHTTPHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /get_all_users", h.GetAllUsers)
	mux.HandleFunc("GET /get_user", h.GetUserByID)

	mux.HandleFunc("GET /get_all_machines", h.GetAllMachines)
	mux.HandleFunc("GET /get_machine", h.GetMachineByID)

	mux.HandleFunc("GET /get_all_sessions", h.GetAllSessions)
	mux.HandleFunc("GET /get_session", h.GetSessionByID)

	// logging all request with LoggingMiddleware
	return middlewares.LoggingMiddleware(mux)
}

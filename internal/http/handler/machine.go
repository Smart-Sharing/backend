package handler

import (
	"log/slog"
	"net/http"

	"github.com/ecol-master/sharing-wh-machines/internal/utils"
)

func (h *Handler) GetAllMachines(w http.ResponseWriter, r *http.Request) {
	machines, err := h.service.GetAllMachines()
	if err != nil {
		slog.Error(
			"failed get all users",
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("error", err.Error()),
		)

		if err := utils.RespondWith500(w); err != nil {
			slog.Error("failed respond with error", slog.Int("status", 500))
		}
		return
	}

	if err := utils.RespondWithJSON(w, 200, machines); err != nil {
		slog.Error("failed to respond with json with users",
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("error", err.Error()),
		)

		if err := utils.RespondWith500(w); err != nil {
			slog.Error("failed respond with error", slog.Int("status", 500))
		}
	}
}

func (h *Handler) GetMachineByID(w http.ResponseWriter, r *http.Request) {
	machineId := r.URL.Query().Get("machine_id")
	machine, err := h.service.GetMachineByID(machineId)

	if err != nil {
		slog.Error("failed to get machine from db",
			slog.String("machine_id", machineId),
			slog.String("error", err.Error()),
		)
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond with 500",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)

		}
		return
	}

	if err = utils.RespondWithJSON(w, 200, machine); err != nil {
		slog.Error("failed to respond with json with machine",
			slog.String("machine_id", machineId),
			slog.String("error", err.Error()),
		)
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond with 500",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
	}
}

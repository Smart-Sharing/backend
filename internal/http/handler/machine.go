package handler

import (
	"log/slog"
	"net/http"

	"github.com/ecol-master/sharing-wh-machines/internal/utils"
)

func (h *Handler) GetAllMachines(w http.ResponseWriter, r *http.Request) {
	op := slog.String("op", "handler.GetAllMachines")

	machines, err := h.service.GetAllMachines()
	if err != nil {
		slog.Error("get all machines", op, slog.String("error", err.Error()))

		if err := utils.RespondWith400(w, "failed to get all machines"); err != nil {
			slog.Error("failed respond with 400", op, slog.String("error", err.Error()))
		}
		return
	}

	if err := utils.RespondWithJSON(w, 200, machines); err != nil {
		slog.Error("failed respond with JSON", op, slog.String("error", err.Error()))
	}
}

func (h *Handler) GetMachineByID(w http.ResponseWriter, r *http.Request) {
	op := slog.String("op", "handler.GetmachineByID")

	machineId := r.URL.Query().Get("machine_id")
	machine, err := h.service.GetMachineByID(machineId)

	if err != nil {
		slog.Error("get machine from db", op, slog.String("machine_id", machineId),
			slog.String("error", err.Error()))

		if err = utils.RespondWith400(w, "failed get machine by id"); err != nil {
			slog.Error("failed to respond with 400", slog.String("error", err.Error()))

		}
		return
	}

	if err = utils.RespondWithJSON(w, 200, machine); err != nil {
		slog.Error("failed to respond with json with machine", op, slog.String("machine_id", machineId),
			slog.String("error", err.Error()))

		if err = utils.RespondWith400(w, "failed respond with JSON"); err != nil {
			slog.Error("failed to respond with 400", slog.String("error", err.Error()))

		}
	}
}

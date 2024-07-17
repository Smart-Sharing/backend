package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/ecol-master/sharing-wh-machines/internal/utils"
)

func (h *Handler) RegisterMachine(w http.ResponseWriter, r *http.Request) {
	var respData struct {
		MachineId string `json:"machine_id"`
		IPAddr    string `json:"ip_addr"`
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed io.ReadAll(r.Body) data on RegisterMachine",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	err = json.Unmarshal(bodyBytes, &respData)
	if err != nil {
		if err = utils.RespondWith400(w, "request body is uncorrect"); err != nil {
			if err = utils.RespondWith500(w); err != nil {
				slog.Error("failed to unmarshal request data on RegisterMachine",
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("error", err.Error()),
				)
			}
		}
		return
	}

	machine, err := h.service.GetMachineByID(respData.MachineId)
	if err != nil {
		// creating new machine if doesn't exists
		machine, err = h.service.InsertMachine(respData.MachineId, respData.IPAddr)
		if err != nil {
			slog.Error("failed to create new in machine",
				slog.String("machineId", respData.MachineId),
				slog.String("ipAddr", respData.IPAddr),
				slog.String("error", err.Error()),
			)

			if err = utils.RespondWith400(w, "failed to create new machine"); err != nil {
				slog.Error("failed to respond with 400 on create new in machine",
					slog.String("machineId", respData.MachineId),
					slog.String("ipAddr", respData.IPAddr),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("error", err.Error()),
				)
			}
			return
		}

	} else {
		machine, err = h.service.UpdateMachineIPAddr(machine.Id, respData.IPAddr)
		if err != nil {
			slog.Error("failed to update machine IP Address",
				slog.String("machineId", respData.MachineId),
				slog.String("ipAddr", respData.IPAddr),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("error", err.Error()),
			)

			if err = utils.RespondWith400(w, "failed to update machine IP Address"); err != nil {
				slog.Error("failed to respond with 400 on failed update machine IP Addr",
					slog.String("machineId", respData.MachineId),
					slog.String("ipAddr", respData.IPAddr),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("error", err.Error()),
				)
			}
			return
		}
	}

	payload := struct {
		CurrentStatus int `json:"current_status"`
	}{CurrentStatus: machine.State}

	if err = utils.SuccessRespondWith200(w, payload); err != nil {
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond Success(200) with paylod on RegisterMachine",
				slog.Any("payload", payload),
				slog.String("machineId", respData.MachineId),
				slog.String("ipAddr", respData.IPAddr),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("error", err.Error()),
			)
		}
	}
}

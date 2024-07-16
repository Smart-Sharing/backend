package handler

import (
	"encoding/json"
	"fmt"
	"io"
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

func (h *Handler) UnlockMachine(w http.ResponseWriter, r *http.Request) {
	// parse MachineId from request
	var respData struct {
		MachineId string `json:"machine_id"`
	}

	dataBytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("failed to read UnlockMachine request body")
		if err := utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond with 500 during failed to unmarshal UnlockMachine request body",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	err = json.Unmarshal(dataBytes, &respData)
	if err != nil {
		slog.Error("failed to unmarshal UnlockMachine request body")
		if err := utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond with 500 during failed to unmarshal UnlockMachine request body",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	machine, err := h.service.GetMachineByID(respData.MachineId)
	if err != nil {
		if err = utils.RespondWith400(w, "machine with such id doesn't exists"); err != nil {
			if err = utils.RespondWith500(w); err != nil {
				slog.Error("failed to respond with 500 during machine with such id doesn't exists",
					slog.String("machine_id", respData.MachineId),
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
					slog.String("error", err.Error()),
				)
			}
		}
		return
	}

	machineSessions, err := h.service.GetActiveSessionsByMachineID(machine.Id)
	if err != nil {
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond with 500 during getting active sessions by machine id",
				slog.String("machine_id", respData.MachineId),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	if len(machineSessions) > 0 {
		msg := fmt.Sprintf("machine with id=%s in using now", respData.MachineId)
		if err = utils.RespondWith400(w, msg); err != nil {
			if err = utils.RespondWith500(w); err != nil {
				slog.Error("failed to respond with 500 during machine in using now",
					slog.String("machine_id", respData.MachineId),
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
					slog.String("error", err.Error()),
				)
			}
		}

		return
	}

	userId, ok := r.Context().Value("user_id").(int)
	if !ok {
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond with 500 on failed get user_id from request context",
				slog.String("machine_id", respData.MachineId),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}

		return
	}

	user, err := h.service.GetUserByID(userId)
	if err != nil {
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond with 500 on failed get user by id",
				slog.String("machine_id", respData.MachineId),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	userSessions, err := h.service.GetActiveSessionsByUserID(user.Id)
	if err != nil {
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond 500 on failed get active sessions by user_id",
				slog.Int("user_id", user.Id),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	if len(userSessions) > 0 {
		if err = utils.RespondWith400(w, "user has another active sessions"); err != nil {
			if err = utils.RespondWith500(w); err != nil {
				slog.Error("failed to respond 500 on user has another active sessions",
					slog.Int("user_id", user.Id),
					slog.Any("user_active_sessions", userSessions),
					slog.String("path", r.URL.Path),
					slog.String("method", r.Method),
					slog.String("error", err.Error()),
				)
			}
		}
		return
	}

	// ВАЖНО Перед созданием сессии надо проверить - может ли этот пользователь взять эту конкретную машину

	// TODO: check machine can be used
	// сделать запрос к ардуинке и проверить что машина активна
	// сделать запрос к ардуино - отправить текущий статус машины
	// ...

	// TODO: check user can rent this machine
	// ...
	// ...
	// ...

	session, err := h.service.InsertSession(user.Id, machine.Id)
	if err != nil {
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond 500 on failed respond with json (session_id)",
				slog.Int("user_id", user.Id),
				slog.String("machine_id", machine.Id),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	payload := struct {
		SessionId int `json:"sessionId"`
	}{SessionId: session.Id}

	if err = utils.SuccessRespondWith200(w, payload); err != nil {
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond 500 on failed respond with json (session_id)",
				slog.Any("payload", payload),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
	}
}

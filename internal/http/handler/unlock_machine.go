package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/ecol-master/sharing-wh-machines/internal/utils"
)

func (h *Handler) UnlockMachine(w http.ResponseWriter, r *http.Request) {
	op := slog.String("op", "handler.UnlockMachine")

	var respData struct {
		MachineId string `json:"machine_id"`
	}

	err := utils.ParseRequestData(r.Body, &respData)
	if err != nil {
		slog.Error("failed parse request data", op, slog.String("error", err.Error()))

		if err := utils.RespondWith500(w); err != nil {
			slog.Error("respond with 500 during failed parse request body",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return

	}

	/*
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
	*/

	machine, err := h.service.GetMachineByID(respData.MachineId)
	if err != nil {
		slog.Error("machine with such id doesn't exists", slog.String("machine_id", respData.MachineId))
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
		slog.Error("failed to get active sessions by machineId",
			slog.String("machine_id", respData.MachineId),
			slog.String("error", err.Error()),
		)
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

	userId, ok := r.Context().Value("user_id").(int64)
	if !ok {
		slog.Error("failed to get user_id from r.Context", slog.Bool("ok", ok))

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

	user, err := h.service.GetUserByID(int(userId))
	if err != nil {
		slog.Error("failed get user by id", slog.Int("user_id", int(userId)), slog.String("error", err.Error()))

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

	if err = canUnlockMachine(h.service, user, machine); err != nil {
		slog.Error("try unlock machine", op, slog.Int("user_id", int(userId)),
			slog.String("machine_id", machine.Id), slog.String("error", err.Error()))

		if err = utils.RespondWith400(w, "user can not unlock machine"); err != nil {
			slog.Error("failed to respond 400 on failed get active sessions by user_id",
				slog.Int("user_id", user.Id),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	machine.State = entities.MachineInUse
	if err = sendMachineCurrentState(machine, h.cfg.MC.RequestTimeout); err != nil {
		slog.Error("failed sendMachineCurrentState", op, slog.String("error", err.Error()))

		if err = utils.RespondWith400(w, "machine can not be used at the current moment"); err != nil {
			slog.Error("failed to respond with 400 on machine is not active",
				slog.Any("machine", machine),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	_, err = h.service.UpdateMachineState(machine.Id, machine.State)
	if err != nil {
		// TODO: подумать, что должно произойти, если не удалось обновить машину
		slog.Error("failed to update machine state UnlockMachine",
			slog.Any("machine", machine),
			slog.Int("new_state", machine.State),
			slog.String("error", err.Error()),
		)

		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond 500 on failed update machine state UnlockMachine",
				slog.String("machine_id", machine.Id),
				slog.Int("new_state", machine.State),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	session, err := h.service.InsertSession(user.Id, machine.Id)
	if err != nil {
		slog.Error("failed to insert new session",
			slog.Int("user_id", int(userId)),
			slog.String("machine_id", machine.Id),
			slog.String("error", err.Error()),
		)
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

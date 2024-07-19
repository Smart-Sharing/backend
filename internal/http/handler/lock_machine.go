package handler

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/ecol-master/sharing-wh-machines/internal/utils"
)

func (h *Handler) LockMachine(w http.ResponseWriter, r *http.Request) {
	// TODO: распарсить данные для работы
	var data struct {
		MachineId string `json:"machine_id"`
	}

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("failed to read r.Body in LockMachine", slog.String("error", err.Error()))
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond 500 on failed to read r.Body in LockMachine",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		slog.Error("failed to unmarhsal bytes from r.Body in LockMachine", slog.String("error", err.Error()))
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond 500 on failed to unmarshal bytes in LockMachine",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	machine, err := h.service.GetMachineByID(data.MachineId)
	if err != nil {
		slog.Error("failed to get machine by id from request data",
			slog.String("data_machine_id", data.MachineId),
			slog.String("error", err.Error()),
		)
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond 500 on failed to get machine by id",
				slog.String("data_machine_id", data.MachineId),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	userId, ok := r.Context().Value("user_id").(int64)
	if !ok {
		slog.Error("failed to get `user_id` value from r.Context", slog.Any("context", r.Context()))

		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond 500 on failed to get 'use_id' value from r.Context",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}

		return
	}

	user, err := h.service.GetUserByID(int(userId))
	if err != nil {
		slog.Error("failed to get user by id", slog.Int64("user_id", userId),
			slog.String("error", err.Error()))

		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond 500 on failed to get user by id",
				slog.Int64("user_id", userId),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	// TODO: проверить, есть ли вообще такая активная сессия (на всякий случай)
	sessions, err := h.service.GetActiveSessionsByMachineAndUser(machine.Id, user.Id)
	if err != nil {
		slog.Error("failed go get sessions by machine and user in LockMachine",
			slog.String("machine_id", machine.Id), slog.Int("user_id", user.Id))

		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond 500 on failed to get user by id",
				slog.Int64("user_id", userId),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}

		return
	}
	if len(sessions) == 0 {
		slog.Error("user has no active session with that machine",
			slog.String("machine_id", machine.Id), slog.Int("user_id", user.Id))

		if err = utils.RespondWith400(w, "user has no active sessions with that machine"); err != nil {
			slog.Error("failed to respond with 400 on user has no active sessions with machine",
				slog.String("machine_id", machine.Id),
				slog.Int("user_id", user.Id),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	if len(sessions) > 1 {
		slog.Error("user has several sessions with a machine", slog.Int("sessions_count", len(sessions)),
			slog.String("machine_id", machine.Id), slog.Int("user_id", user.Id))

		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond 500 on failed to get user by id",
				slog.Int64("user_id", userId),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	// TODO: обновить данные в базе данных у машины
	machine.State = entities.MachineFree
	_, err = h.service.UpdateMachineState(machine.Id, machine.State)
	if err != nil {
		slog.Error("failed to update machine state LockMachine",
			slog.Any("machine", machine),
			slog.Int("new_state", machine.State),
			slog.String("error", err.Error()),
		)

		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond 500 on failed update machine state LockMachine",
				slog.String("machine_id", machine.Id),
				slog.Int("new_state", machine.State),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}

		return
	}
	// TODO: Отправить данные на машину, для обработки
	if recieved := sendMachineCurrentState(machine, h.cfg.MC.RequestTimeout); !recieved {
		slog.Error("failed to send machine new current state Lock.Machine",
			slog.String("machine_id", machine.Id), slog.Int("new_state", machine.State))
	}

	// TODO: завершить сессию
	session, err := h.service.UpdateSessionState(sessions[0].Id, entities.SessionStopped)
	if err != nil {
		// TODO:
		slog.Error("failed to update session state UnlockMachine",
			slog.Any("machine", machine),
			slog.Int("new_state", machine.State),
			slog.String("error", err.Error()),
		)

		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond 500 on failed session machine state UnlockMachine",
				slog.Int("session_id", session.Id),
				slog.Int("new_state", entities.SessionStopped),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}

		return
	}

	// TODO: ответить, что добавление сессии прошло успешно
	slog.Info("session was successfully stopped", slog.Int("session_id", session.Id))

	if err = utils.SuccessRespondWith200(w, "successfullly lock machine"); err != nil {
		slog.Error("failed to respond with 200 on lock machine",
			slog.String("machine_id", machine.Id),
			slog.Int("user_id", user.Id),
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("error", err.Error()),
		)
	}
}

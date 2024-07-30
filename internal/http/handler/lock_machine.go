package handler

import (
	"log/slog"
	"net/http"

	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/ecol-master/sharing-wh-machines/internal/utils"
)

func (h *Handler) LockMachine(w http.ResponseWriter, r *http.Request) {
	op := slog.String("op", "handler.LockMachine")
	// TODO: распарсить данные для работы
	var data struct {
		MachineId string `json:"machine_id"`
	}

	err := utils.ParseRequestData(r.Body, &data)
	if err != nil {
		slog.Error("failed parse request data", op, slog.String("error", err.Error()))
		if err = utils.RespondWith400(w, "failed parse request body"); err != nil {
			slog.Error("failed respond with 400", slog.String("error", err.Error()))
		}
		return
	}

	machine, err := h.service.GetMachineByID(data.MachineId)
	if err != nil {
		slog.Error("get machine by id", op, slog.String("machine_id", data.MachineId),
			slog.String("error", err.Error()))

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
		slog.Error("get `user_id` from r.Context", op, slog.Any("context", r.Context()))

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
		slog.Error("get user by id", op, slog.Int64("user_id", userId), slog.String("error", err.Error()))

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

	session, err := canLockMachine(h.service, user, machine)
	if err != nil {
		slog.Error("tryLockMachine", op, slog.Int("user_id", user.Id),
			slog.String("machine_id", machine.Id), slog.String("error", err.Error()))

		if err = utils.RespondWith400(w, "user can not lock machine"); err != nil {
			slog.Error("failed to respond 500 on failed to try lock machine",
				slog.Int64("user_id", userId),
				slog.String("machine_id", machine.Id),
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
	if err := sendMachineCurrentState(machine, h.cfg.MC.RequestTimeout); err != nil {
		slog.Error("send machine new state", slog.String("machine_id", machine.Id),
			slog.Int("new_state", machine.State), slog.String("error", err.Error()))
	}

	// TODO: завершить сессию
	session, err = h.service.UpdateSessionState(session.Id, entities.SessionStopped)
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

	payload := struct {
		Msg string `json:"msg"`
	}{Msg: "successfullly lock machine"}

	if err = utils.SuccessRespondWith200(w, payload); err != nil {
		slog.Error("failed to respond with 200 on lock machine",
			slog.String("machine_id", machine.Id),
			slog.Int("user_id", user.Id),
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("error", err.Error()),
		)
	}
}

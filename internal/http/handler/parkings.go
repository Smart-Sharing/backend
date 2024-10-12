package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/ecol-master/sharing-wh-machines/internal/utils"
)

func (h *Handler) GetAllParkings(w http.ResponseWriter, r *http.Request) {
	parkings, err := h.service.GetAllParkings()

	if err != nil {
		slog.Error(
			"failed get all parkings",
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("error", err.Error()),
		)

		if err := utils.RespondWith500(w); err != nil {
			slog.Error("failed respond with error", slog.Int("status", 500))
		}
		return
	}

	if err := utils.RespondWithJSON(w, 200, parkings); err != nil {
		slog.Error("failed to respond with json with parkings",
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("error", err.Error()),
		)

		if err := utils.RespondWith500(w); err != nil {
			slog.Error("failed respond with error", slog.Int("status", 500))
		}
	}
}

func (h *Handler) GetParkingById(w http.ResponseWriter, r *http.Request) {
	parkingId := r.URL.Query().Get("parking_id")

	id, err := strconv.Atoi(parkingId)
	if err != nil {
		slog.Error("`parking_id` query is not integer")
		if err := utils.RespondWith400(w, "parking_id should be integer"); err != nil {
			slog.Error(
				"failed respond 400",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	parking, err := h.service.GetParkingById(id)
	if err != nil {
		slog.Error("parking not found", slog.Int("parking_id", id))
		if err := utils.RespondWith400(w, "parking not found"); err != nil {
			slog.Error(
				"failed respond 400",
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	if err := utils.RespondWithJSON(w, 200, parking); err != nil {
		slog.Error("failed to respond with json with parking",
			slog.Int("parking_id", id),
			slog.String("error", err.Error()),
		)
		if err := utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond with 500 during get parking by id",
				slog.Int("parking_id", id),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
	}
}

func (h *Handler) RegisterParking(w http.ResponseWriter, r *http.Request) {
	op := slog.String("op", "handler.RegisterParking")

	data := entities.Parking{}
	if err := utils.ParseRequestData(r.Body, &data); err != nil {
		slog.Error("parse req data", op, slog.String("error", err.Error()))
		if err := utils.RespondWith400(w, "failed to parse request data"); err != nil {
			slog.Error("failed to respond with 400", op, slog.String("error", err.Error()))
		}
		return
	}

	name := slog.String("parkingName", data.Name)
	mac := slog.String("machineId", data.MacAddr)
	cap := slog.Int("machineId", int(data.Capacity))
	state := slog.Int("machineId", int(data.State))

	parking, err := h.service.InsertParking(data.Name, data.MacAddr, data.Capacity, data.State)
	if err != nil {
		slog.Error("failed to create new in parking. Maybe, parking with this name already exists", name, mac, cap, state, slog.String("error", err.Error()))
		if err = utils.RespondWith400(w, "failed to create new in parking. Maybe, parking with this name already exists"); err != nil {
			slog.Error("failed to respond with 400", op, slog.String("error", err.Error()))
		}
		return
	}

	if err = utils.SuccessRespondWith200(w, parking); err != nil {
		if err = utils.RespondWith500(w); err != nil {
			slog.Error("failed to respond Success(200) with paylod on RegisterMachine",
				slog.Any("payload", parking), name, mac, cap, state,
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("error", err.Error()),
			)
		}
	}
}

func (h *Handler) UpdateParkingState(w http.ResponseWriter, r *http.Request) {
	op := slog.String("op", "handler.UpdateParking")

	data := struct {
		ParkingId int `json:"id"`
		NewState  int `json:"state"`
	}{}

	err := utils.ParseRequestData(r.Body, &data)
	if err != nil {
		slog.Error("failed parse request data", op, slog.String("error", err.Error()))
		if err = utils.RespondWith400(w, "failed parse request body"); err != nil {
			slog.Error("failed respond with 400", slog.String("error", err.Error()))
		}
		return
	}

	if data.NewState > 1 || data.NewState < 0 {
		if err = utils.RespondWith400(w, "error while updating parking state. Invalid parking state. Use 0 or 1"); err != nil {
			slog.Error("failed to respond 400 on failed update parking state",
				slog.Int("parking_id", data.ParkingId),
				slog.Int("new_state", int(data.NewState)),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	parking, err := h.service.UpdateParkingState(entities.ParkingState(data.NewState), data.ParkingId)
	if err != nil {
		slog.Error("failed to update parking state",
			slog.Any("parking", data),
			slog.Int("new_state", data.NewState),
			slog.String("error", err.Error()),
		)

		if err = utils.RespondWith400(w, "error while updating parking state. Parking not exists or missing field id"); err != nil {
			slog.Error("failed to respond 400 on failed update parking state",
				slog.Int("parking_id", data.ParkingId),
				slog.Int("new_state", data.NewState),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	if err = utils.SuccessRespondWith200(w, parking); err != nil {
		slog.Error("failed to respond with 200 on lock machine",
			slog.Int("parking_id", data.ParkingId),
			slog.Int("new_state", data.NewState),
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("error", err.Error()),
		)
	}
}

func (h *Handler) UpdateParkingCapacity(w http.ResponseWriter, r *http.Request) {
	op := slog.String("op", "handler.ManualyAddParkingMachine")

	data := struct {
		ParkingId   int               `json:"id"`
		NewCapacity entities.Capacity `json:"capacity"`
	}{}

	err := utils.ParseRequestData(r.Body, &data)
	if err != nil {
		slog.Error("failed parse request data", op, slog.String("error", err.Error()))
		if err = utils.RespondWith400(w, "failed parse request body"); err != nil {
			slog.Error("failed respond with 400", slog.String("error", err.Error()))
		}
		return
	}

	parking, err := h.service.UpdateParkingCapacity(data.NewCapacity, data.ParkingId)
	if err != nil {
		slog.Error("failed to update parking capacity",
			slog.Int("parking_id", data.ParkingId),
			slog.String("error", err.Error()),
		)

		if err = utils.RespondWith400(w, "error while updating parking capacity. Parking not exists or missing field id"); err != nil {
			slog.Error("error while updating parking capacity",
				slog.Any("parking", data),
				slog.Int("new_machines", parking.Machines+1),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	if err = utils.SuccessRespondWith200(w, parking); err != nil {
		slog.Error("failed to respond with 200 on update machine parking id",
			slog.Any("parking", parking),
			slog.Int("parking_id", data.ParkingId),
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("error", err.Error()),
		)
	}
}

func (h *Handler) ManualyAddParkingMachine(w http.ResponseWriter, r *http.Request) {
	op := slog.String("op", "handler.ManualyAddParkingMachine")

	data := struct {
		MachineId string `json:"machine_id"`
		ParkingId int    `json:"parking_id"`
	}{}

	err := utils.ParseRequestData(r.Body, &data)
	if err != nil {
		slog.Error("failed parse request data", op, slog.String("error", err.Error()))
		if err = utils.RespondWith400(w, "failed parse request body"); err != nil {
			slog.Error("failed respond with 400", slog.String("error", err.Error()))
		}
		return
	}

	parking, err := h.service.GetParkingById(data.ParkingId)
	if err != nil {
		slog.Error("failed to get parking by id",
			slog.Int("parkingId", data.ParkingId),
			slog.String("error", err.Error()),
		)

		if err = utils.RespondWith400(w, "error while updating parking machines. Parking not exists or missing field parking_id"); err != nil {
			slog.Error("failed to respond 400 on failed update parking state",
				slog.Int("parking_id", data.ParkingId),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	if int(parking.Capacity) <= parking.Machines && parking.Capacity != 0 {
		if err = utils.RespondWith400(w, "error while updating parking state. Parking machines is more or equals than capacity"); err != nil {
			slog.Error("failed to respond 400 on failed update parking state",
				slog.Int("parking_id", data.ParkingId),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	parking, err = h.service.UpdateParkingMachines(parking.Machines+1, parking.Id)
	if err != nil {
		slog.Error("failed to update parking machines",
			slog.Any("parking", data),
			slog.Int("new_machines", parking.Machines+1),
			slog.String("error", err.Error()),
		)

		if err = utils.RespondWith400(w, "error while updating parking machines"); err != nil {
			slog.Error("error while updating parking machines",
				slog.Any("parking", data),
				slog.Int("new_machines", parking.Machines+1),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	_, err = h.service.UpdateMachineParkingId(data.MachineId, data.ParkingId)
	if err != nil {
		slog.Error("failed to update parking_id for machine",
			slog.String("machineId", data.MachineId),
			slog.Int("parkingId", data.ParkingId),
			slog.String("error", err.Error()),
		)

		if err = utils.RespondWith400(w, "error while updating parking_id for machine. Machine not exists or missing field machine_id"); err != nil {
			slog.Error("error while get parking by id",
				slog.String("machineId", data.MachineId),
				slog.Int("parkingId", data.ParkingId),
				slog.String("path", r.URL.Path),
				slog.String("method", r.Method),
				slog.String("error", err.Error()),
			)
		}
		return
	}

	if err = utils.SuccessRespondWith200(w, parking); err != nil {
		slog.Error("failed to respond with 200 on update machine parking id",
			slog.String("machine_id", data.MachineId),
			slog.Int("parking_id", data.ParkingId),
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("error", err.Error()),
		)
	}
}

package handler

import (
	"encoding/json"
	"io"
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
		return
	}

	err = json.Unmarshal(bodyBytes, &respData)
	if err != nil {
		return
	}

	machine, err := h.service.GetMachineByID(respData.MachineId)
	if err != nil {
		machine, err = h.service.InsertMachine(respData.MachineId, respData.IPAddr)
		if err != nil {
			// TODO: failed to crete new machine
			return
		}

	} else {
		machine, err = h.service.UpdateMachineIPAddr(machine.Id, respData.IPAddr)
		if err != nil {
			// TODO: process error if can not update machine ip addr
		}
	}

	payload := struct {
		CurrentStatus int `json:"current_status"`
	}{CurrentStatus: machine.State}
	if err = utils.SuccessRespondWith200(w, payload); err != nil {
		if err = utils.RespondWith500(w); err != nil {
			// TODO: add failed logs
		}
	}
}

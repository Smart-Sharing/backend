package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ecol-master/sharing-wh-machines/internal/entities"
)

func sendMachineCurrentState(machine *entities.Machine, timeout time.Duration) bool {
	payload := []byte(fmt.Sprintf(`{"current_state": %d}`, machine.State))
	reader := bytes.NewReader(payload)

	address := fmt.Sprintf("http://%s/%s", machine.IPAddr, machine.Id)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, address, reader)
	if err != nil {
		// TODO: что-то придумать, когда не получается отправить запрос к ардуино
		return false
	}
	resp, err := (&http.Client{}).Do(req)

	if err != nil {
		// TODO: что то придумать для обработки
		return false
	}

	return resp.StatusCode == http.StatusOK
}

package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/ecol-master/sharing-wh-machines/internal/service"
	"github.com/pkg/errors"
)

func sendMachineCurrentState(machine *entities.Machine, timeout time.Duration) error {
	payload := []byte(fmt.Sprintf(`{"current_state": %d}`, machine.State))
	reader := bytes.NewReader(payload)

	address := fmt.Sprintf("http://%s/%s", machine.IPAddr, machine.Id)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, address, reader)
	if err != nil {
		return errors.Wrap(err, "create new request")
	}
	resp, err := (&http.Client{}).Do(req)

	if err != nil {
		return errors.Wrap(err, "do request")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed send to arduino current status")
	}

	return nil
}

func canUnlockMachine(svc *service.Service, user *entities.User, _ *entities.Machine) error {
	switch user.JobPosition {
	case entities.Worker:
		userSessions, err := svc.GetActiveSessionsByUserID(user.Id)
		if err != nil {
			return errors.Wrap(err, "get active sessions by userId")
		}

		if len(userSessions) > 0 {
			msg := fmt.Sprintf("user has several active sessions, cnt=%d", len(userSessions))
			return errors.Wrap(err, msg)
		}
		return nil

	case entities.Admin:
		return nil

	default:
		return errors.New("user has uknown job position")
	}
}

func canLockMachine(svc *service.Service, user *entities.User, machine *entities.Machine) (*entities.Session, error) {
	if user.JobPosition == entities.Worker {
		sessions, err := svc.GetActiveSessionsByMachineAndUser(machine.Id, user.Id)
		if err != nil {
			return nil, errors.Wrap(err, "get sessions by machine.Id and user.Id")
		}

		if len(sessions) == 0 {
			return nil, errors.New("user has no active sessions with that machine")
		}

		if len(sessions) > 1 {
			return nil, errors.New("user has several active sessions with machine")
		}

		return &sessions[0], nil
	}

	if user.JobPosition == entities.Admin {
		sessions, err := svc.GetActiveSessionsByMachineID(machine.Id)
		if err != nil {
			return nil, errors.Wrap(err, "get active sessions by machine.Id")
		}
		if len(sessions) == 0 {
			return nil, errors.New("there is no active sessions with machine")
		}
		if len(sessions) > 1 {
			return nil, errors.New("there several active sessions with machine")
		}
		return &sessions[0], nil
	}

	return nil, errors.New("user has uknown job position")

}

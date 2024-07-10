package service

import (
	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/ecol-master/sharing-wh-machines/internal/repositories/machines"
	"github.com/ecol-master/sharing-wh-machines/internal/repositories/sessions"
	"github.com/ecol-master/sharing-wh-machines/internal/repositories/users"
	"github.com/jmoiron/sqlx"
)

type User interface {
	GetUserByID(userId int) (*entities.User, error)
	GetAllUsers() ([]entities.User, error)
}

type Machine interface {
	GetMachineByID(machineId string) (*entities.Machine, error)
	GetAllMachines() ([]entities.Machine, error)
}

type Session interface {
	GetSessionByID(sessionId int) (*entities.Session, error)
	GetAllSessions() ([]entities.Session, error)
}

type Service struct {
	User
	Machine
	Session
}

func New(db *sqlx.DB) *Service {
	return &Service{
		User:    users.NewRepository(db),
		Machine: machines.NewRepository(db),
		Session: sessions.NewRepository(db),
	}
}

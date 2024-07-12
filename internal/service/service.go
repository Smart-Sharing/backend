package service

import (
	"time"

	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/ecol-master/sharing-wh-machines/internal/jwt"
	"github.com/ecol-master/sharing-wh-machines/internal/repositories/machines"
	"github.com/ecol-master/sharing-wh-machines/internal/repositories/sessions"
	"github.com/ecol-master/sharing-wh-machines/internal/repositories/users"
	"github.com/jmoiron/sqlx"
)

type User interface {
	GetAllUsers() ([]entities.User, error)
	GetUserByID(userId int) (*entities.User, error)
	GetUserByPhoneNumber(phoneNumber string) (*entities.User, error)
}

type Machine interface {
	GetMachineByID(machineId string) (*entities.Machine, error)
	GetAllMachines() ([]entities.Machine, error)
}

type Session interface {
	GetSessionByID(sessionId int) (*entities.Session, error)
	GetAllSessions() ([]entities.Session, error)
}

type Auth interface {
	GenerateToken(user entities.User, secret string, tokenTTL time.Duration) (string, error)
}

type Service struct {
	User
	Machine
	Session
	Auth
}

func New(db *sqlx.DB) *Service {
	return &Service{
		User:    users.NewRepository(db),
		Machine: machines.NewRepository(db),
		Session: sessions.NewRepository(db),
		Auth:    jwt.NewService(),
	}
}

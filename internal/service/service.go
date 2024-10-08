package service

import (
	"time"

	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/ecol-master/sharing-wh-machines/internal/libs/jwt"
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
	InsertMachine(machineId, ipAddr string) (*entities.Machine, error)
	GetMachineByID(machineId string) (*entities.Machine, error)
	GetAllMachines() ([]entities.Machine, error)
	UpdateMachineIPAddr(machineId, ipAddr string) (*entities.Machine, error)
	UpdateMachineState(machineId string, state entities.MachineState) (*entities.Machine, error)
}

type Session interface {
	InsertSession(workerId int, machineId string) (*entities.Session, error)
	GetSessionByID(sessionId int) (*entities.Session, error)
	GetAllSessions() ([]entities.Session, error)

	GetActiveSessionsByMachineID(machineId string) ([]entities.Session, error)
	GetPausedSessionsByMachineID(machineId string) ([]entities.Session, error)

	GetActiveSessionsByUserID(userId int) ([]entities.Session, error)
	GetPauseSessionsByUserID(userId int) ([]entities.Session, error)
	GetUnfinishedSessionsByUserId(userId int) ([]entities.Session, error)

	GetActiveSessionsByMachineAndUser(machineId string, userId int) ([]entities.Session, error)
	GetPausedSessionsByMachineAndUser(machineId string, userId int) ([]entities.Session, error)

	UpdateSessionState(sessionId int, state entities.SessionState) (*entities.Session, error)

	PauseSession(sessionId int) (*entities.Session, error)
	FinishSession(sessionId int) (*entities.Session, error)
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

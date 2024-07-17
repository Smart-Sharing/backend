package sessions

import (
	"github.com/ecol-master/sharing-wh-machines/internal/entities"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) InsertSession(userId int, machineId string) (*entities.Session, error) {
	var (
		session entities.Session
		id      int
	)

	q := `INSERT INTO sessions (machine_id, worker_id) VALUES ($1, $2) RETURNING id;`

	if err := r.db.QueryRowx(q, machineId, userId).Scan(&id); err != nil {
		return nil, errors.Wrap(err, "insert new sessiona and scan id")
	}

	q = `SELECT * FROM sessions WHERE id = $1`
	if err := r.db.Get(&session, q, id); err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *repository) GetSessionByID(sessionId int) (*entities.Session, error) {
	var session entities.Session

	q := `SELECT * FROM sessions WHERE id = $1`
	if err := r.db.Get(&session, q, sessionId); err != nil {
		return nil, errors.Wrap(err, "get session by id")
	}
	return &session, nil
}

func (r *repository) GetAllSessions() ([]entities.Session, error) {
	sessions := make([]entities.Session, 0)

	q := `SELECT * FROM sessions`
	if err := r.db.Select(&sessions, q); err != nil {
		return nil, errors.Wrap(err, "select all sessions")
	}
	return sessions, nil
}

func (r *repository) GetActiveSessionsByMachineID(machineId string) ([]entities.Session, error) {
	sessions := make([]entities.Session, 0)

	q := `SELECT * FROM sessions WHERE machine_id = $1 AND state = 0`
	if err := r.db.Select(&sessions, q, machineId); err != nil {
		return nil, errors.Wrap(err, "select all sessions by machineId")
	}
	return sessions, nil
}

func (r *repository) GetActiveSessionsByUserID(userId int) ([]entities.Session, error) {
	sessions := make([]entities.Session, 0)

	q := `SELECT * FROM sessions WHERE worker_id = $1 AND state = 0`
	if err := r.db.Select(&sessions, q, userId); err != nil {
		return nil, errors.Wrap(err, "select all sessions by userId")
	}
	return sessions, nil
}

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

func (r *repository) InsertSession(workerId int, machineId string) error {
	return nil
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

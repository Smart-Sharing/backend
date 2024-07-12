package machines

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

func (r *repository) InsertMachine(machineId string) error {
	q := `INSERT INTO machines (id) VALUES ($1)`
	if _, err := r.db.Exec(q, machineId); err != nil {
		return errors.Wrap(err, "insert machine")
	}
	return nil
}

func (r *repository) GetMachineByID(machineId string) (*entities.Machine, error) {
	var m entities.Machine

	q := `SELECT * FROM machines WHERE id = $1`

	if err := r.db.Get(&m, q, machineId); err != nil {
		return nil, errors.Wrap(err, "select machine by id")
	}
	return &m, nil
}

func (r *repository) GetAllMachines() ([]entities.Machine, error) {
	machines := make([]entities.Machine, 0)
	q := `SELECT * FROM machines`
	if err := r.db.Select(&machines, q); err != nil {
		return nil, errors.Wrap(err, "get all machines")
	}

	return machines, nil
}

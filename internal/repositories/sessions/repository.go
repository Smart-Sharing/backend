package sessions

import "github.com/jmoiron/sqlx"

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) InsertSession(workerId int, machineId string) error {

	return nil
}

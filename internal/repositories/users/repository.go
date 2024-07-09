package users

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

func (r *repository) GetAllUsers() ([]entities.User, error) {
	users := make([]entities.User, 0)

	q := `SELECT * FROM users`
	err := r.db.Select(&users, q)

	if err != nil {
		return nil, errors.Wrap(err, "select all users")
	}
	return users, nil
}
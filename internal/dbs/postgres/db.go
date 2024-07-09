package postgres

import (
	"fmt"

	"github.com/ecol-master/sharing-wh-machines/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func New(cfg config.PostgresConfig) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable", cfg.User, cfg.Password, cfg.Addr, cfg.Port, cfg.DB)

	conn, err := sqlx.Connect("postgres", dataSource)
	if err != nil {
		return nil, errors.Wrap(err, "sqlx connect")
	}

	err = conn.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "conn ping failed")
	}

	return conn, err
}

package app

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ecol-master/sharing-wh-machines/internal/config"
	"github.com/ecol-master/sharing-wh-machines/internal/dbs/postgres"
	"github.com/ecol-master/sharing-wh-machines/internal/http/handler"
	"github.com/pkg/errors"
)

type App struct {
	server *http.Server
	cfg    *config.Config
}

func New(cfg *config.Config) *App {
	return &App{
		server: &http.Server{},
		cfg:    cfg,
	}
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "new user")
}

func (a *App) Run() error {
	db, err := postgres.New(a.cfg.Postgres)
	if err != nil {
		// TODO: посмотреть как делал панику Николай Тузов
		panic(errors.Wrap(err, ""))
	}
	handler := handler.New(db, a.cfg).MakeHTTPHandler()

	addr := fmt.Sprintf("%s:%d", a.cfg.App.Addr, a.cfg.App.Port)
	slog.Info("staring app", slog.String("address", addr))
	return http.ListenAndServe(addr, handler)
}

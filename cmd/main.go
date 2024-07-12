package main

import (
	"log/slog"

	"github.com/ecol-master/sharing-wh-machines/internal/app"
	"github.com/ecol-master/sharing-wh-machines/internal/config"
	"github.com/ecol-master/sharing-wh-machines/internal/logger"
)

func main() {
	logger.Setup()
	cfg := config.MustLoad()

	a := app.New(cfg)
	if err := a.Run(); err != nil {
		slog.Error(err.Error())
	}
}

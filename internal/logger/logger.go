package logger

import (
	"log/slog"
	"os"
)

func Setup() error {
	file, err := os.OpenFile("info.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	logger := slog.New(slog.NewTextHandler(file, nil))
	slog.SetDefault(logger)
	return nil
}

func Info() {}

func Warning() {}

func Error() {}

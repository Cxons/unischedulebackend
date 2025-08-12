package logging

import (
	"log/slog"
	"os"
)


func NewLogger(enviroment string) *slog.Logger{
	var handler slog.Handler

	handleoptions := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	if enviroment == "PRODUCTION" {
		handler = slog.NewJSONHandler(os.Stderr,handleoptions)
	} else{
		handler = slog.NewTextHandler(os.Stderr,handleoptions);
	}
	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger
}
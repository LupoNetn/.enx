package logger

import (
	"log/slog"
	"os"
)

func InitLogger(env string) {
	var handler slog.Handler

	if env == "production" {
		//JSON logs in production easy ot parse and injest
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		//Human readable logs in development
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

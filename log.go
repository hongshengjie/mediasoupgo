package mediasoupgo

import (
	"log/slog"
	"os"
)

func init() {
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true},
		),
	)
	slog.SetDefault(logger)
}

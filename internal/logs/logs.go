package logs

import (
	"log"
	"log/slog"
	"os"

	"github.com/Aoladiy/go-with-tools/internal/config"
)

func Init(c config.Config) {
	var logLevel slog.Level
	switch c.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		log.Fatal("unknown logLevel type in config")
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))
}

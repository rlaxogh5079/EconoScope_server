package logger

import (
	"io"
	"log"
	"os"

	"github.com/rlaxogh5079/EconoScope/config"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	cfg := config.AppConfig
	Log = logrus.New()

	level, err := logrus.ParseLevel(cfg.Logging.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Log.SetLevel(level)

	if cfg.Logging.Format == "json" {
		Log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	switch cfg.Logging.Output {
	case "stdout":
		Log.SetOutput(os.Stdout)
	case "stderr":
		Log.SetOutput(os.Stderr)
	case "file":
		file, err := os.OpenFile(cfg.Logging.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		Log.SetOutput(io.MultiWriter(file, os.Stdout))
	default:
		Log.SetOutput(os.Stdout)
	}
}
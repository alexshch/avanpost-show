package config

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/ilyakaznacheev/cleanenv"
)

type Loader[T any] struct {
	Env         string
	DefaultConf string
}

func (l Loader[T]) Load() (*T, error) {
	path, _ := os.LookupEnv(l.Env)
	if path == "" {
		slog.Info("env variable doesn't set, use default config path", "env", l.Env)
		path = l.DefaultConf
		binPath, err := os.Executable()
		if err != nil {
			slog.Error("load config error", "err", err)
		}
		binDir := filepath.Dir(binPath)
		path = filepath.Join(binDir, l.DefaultConf)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		slog.Error("config file doesn't exist", "path", path)
		return nil, err
	}
	var config T
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func GetLogLevel(l string) slog.Level {
	switch l {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	}
	return slog.LevelInfo
}

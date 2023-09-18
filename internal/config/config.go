package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	log "log/slog"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/obalunenko/orderpacker/internal/packer"
)

type httpConfig struct {
	Port string `yaml:"port" json:"port"`
	Host string `yaml:"host" json:"host"`
}

type packConfig struct {
	Boxes []uint `yaml:"boxes" json:"boxes"`
}

type logConfig struct {
	Level  string `yaml:"level" json:"level"`
	Format string `yaml:"format" json:"format"`
}

type Config struct {
	HTTP httpConfig `yaml:"http" json:"http"`
	Pack packConfig `yaml:"pack" json:"pack"`
	Log  logConfig  `yaml:"log" json:"log"`
}

func DefaultConfig() *Config {
	return &Config{
		HTTP: httpConfig{
			Port: "8080",
			Host: "0.0.0.0",
		},
		Pack: packConfig{
			Boxes: packer.DefaultBoxes,
		},
		Log: logConfig{
			Level:  "INFO",
			Format: "text",
		},
	}
}

var (
	ErrEmptyPath = errors.New("empty path")
	ErrNotExists = errors.New("config file not found")
)

func Load(path string) (*Config, error) {
	if path == "" {
		return nil, ErrEmptyPath
	}

	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%w: %s", ErrNotExists, path)
		}

		return nil, fmt.Errorf("failed to open config file: %w", err)
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.Error("Error closing config file", "error", err)
		}
	}()

	var cfg Config

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var unmarshalFn func([]byte, interface{}) error

	switch filepath.Ext(path) {
	case ".json":
		unmarshalFn = json.Unmarshal
	case ".yaml", ".yml":
		unmarshalFn = yaml.Unmarshal
	default:
		return nil, fmt.Errorf("unsupported config file extension: %s", filepath.Ext(path))
	}

	if err = unmarshalFn(b, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return &cfg, nil
}

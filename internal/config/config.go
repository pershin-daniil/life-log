// Package config provides application configuration.
package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	// DefaultPort is the default HTTP server port.
	DefaultPort = 8080
	// DefaultConfigPath is the default path to config file.
	DefaultConfigPath = "config.yaml"
	// DefaultShutdownTimeout is the default shutdown timeout.
	DefaultShutdownTimeout = 3 * time.Second
	// DefaultReadHeaderTimeout is the default read header timeout.
	DefaultReadHeaderTimeout = 1 * time.Second
)

// Config represents the application configuration.
type Config struct {
	// Application settings from YAML
	Server ServerConfig `yaml:"server"`

	// Secrets and environment-specific values from ENV
	Secrets SecretsConfig `yaml:"-"`
}

// ServerConfig contains server-related configuration.
type ServerConfig struct {
	Port              int           `yaml:"port"`
	ShutdownTimeout   time.Duration `yaml:"shutdown_timeout_seconds"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout_seconds"`
}

// SecretsConfig contains secrets and environment-specific values.
type SecretsConfig struct {
	VersionTag       string
	DatabaseURL      string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
}

// Load loads the configuration from file and environment variables.
// Environment variables take precedence over the configuration file.
func Load(lg *slog.Logger) (*Config, error) {
	// Default values
	cfg := &Config{
		Server: ServerConfig{
			Port:              DefaultPort,
			ShutdownTimeout:   DefaultShutdownTimeout,
			ReadHeaderTimeout: DefaultReadHeaderTimeout,
		},
	}

	// Loading from configuration file
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = DefaultConfigPath
	}

	// Validate and clean the config path
	configPath = filepath.Clean(configPath)

	// Check if file exists
	if _, err := os.Stat(configPath); err == nil {
		file, err := os.Open(configPath)
		if err != nil {
			return nil, fmt.Errorf("open config file: %v", err)
		}

		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				// We can't do much about the close error except log it
				lg.Error("error closing config file", "error", closeErr)
			}
		}()

		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(cfg); err != nil {
			return nil, fmt.Errorf("decode config file: %v", err)
		}
	}

	// Environment variables take precedence over config file
	if portStr := os.Getenv("SERVER_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			cfg.Server.Port = port
		}
	}

	// Load secrets from environment variables
	cfg.Secrets = SecretsConfig{
		VersionTag:       os.Getenv("VERSION_TAG"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
	}

	return cfg, nil
}

// ServerAddr returns the formatted server address string.
func (c *Config) ServerAddr() string {
	return fmt.Sprintf(":%d", c.Server.Port)
}

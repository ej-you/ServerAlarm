// Package config provides loading config data from external sources
// like env variables, yaml-files etc.
package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	// db
	_defDBUser    = "root"          // default db user
	_defDBName    = "server_room"   // default db name
	_defDBHost    = "127.0.0.1"     // default db host
	_defDBPort    = "3306"          // default db port
	_defDBTimeout = 5 * time.Second // default db timeout

	// ntfy
	_defNtfyHost  = "127.0.0.1" // default ntfy host
	_defNtfyPort  = "80"        // default ntfy port
	_defNtfyTheme = "server"    // default ntfy theme

	// logging
	_defLogLevel   = 2     // default log level (info)
	_defJSONFormat = false // default log JSON-format
)

type (
	Config struct {
		DB      `yaml:"db"`
		Ntfy    `yaml:"ntfy"`
		Logging `yaml:"logging"`
	}

	DB struct {
		User     string        `yaml:"user"`
		Password string        `env-required:"true" env:"DB_PASSWORD"`
		Name     string        `yaml:"name"`
		Host     string        `yaml:"host"`
		Port     string        `yaml:"port"`
		Timeout  time.Duration `yaml:"timeout"`
		DSN      string
	}

	Ntfy struct {
		Host  string `yaml:"host"`
		Port  string `yaml:"port"`
		Theme string `yaml:"theme"`
		Token string `env-required:"true" env:"NTFY_TOKEN"`
	}

	Logging struct {
		LogLevel int `yaml:"log_level"`
		// use JSON format if true
		JSONFormat bool `yaml:"json_format"`
	}
)

// NewDefault returns a new instance of Config with default data.
func NewDefault() *Config {
	return &Config{
		DB: DB{
			User:    _defDBUser,
			Name:    _defDBName,
			Host:    _defDBHost,
			Port:    _defDBPort,
			Timeout: _defDBTimeout,
		},
		Ntfy: Ntfy{
			Host:  _defNtfyHost,
			Port:  _defNtfyPort,
			Theme: _defNtfyTheme,
		},
		Logging: Logging{
			LogLevel:   _defLogLevel,
			JSONFormat: _defJSONFormat,
		},
	}
}

// Returns app config loaded from YAML-file.
func New() (*Config, error) {
	// fill with default values (for settings in yaml-file)
	cfg := NewDefault()

	// read YAML config file
	if err := cleanenv.ReadConfig("./config.yml", cfg); err != nil {
		return nil, fmt.Errorf("read yaml config file: %w", err)
	}
	// read ENV variables
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("read env-variables: %w", err)
	}

	// collect DSN string
	cfg.DB.DSN = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&timeout=%s",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.Timeout,
	)
	return cfg, nil
}

// Package config provides loading config data from external sources
// like env variables, yaml-files etc.
package config

import (
	"fmt"
	"time"
	_ "time/tzdata" // compressed tz data to don't depend on system files

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

	// app
	_defLocaleName = "UTC" // default locale name
	_defLogLevel   = 2     // default log level (info)
	_defJSONFormat = false // default log JSON-format

	// healthcheck
	_defHealthcheckPort = "80" // default port for healthcheck service

	// climate
	_defCheckDBEvery         = 5 * time.Minute // default time between db data checks
	_defTempTresholdStandart = 35.0            // default standart temperature threshold
	_defTempTresholdHigh     = 40.0            // default high temperature threshold
	_defTempTresholdUrgent   = 45.0            // default urgent temperature threshold
)

type (
	Config struct {
		DB   `yaml:"db"`
		Ntfy `yaml:"ntfy"`
		App  `yaml:"app"`
	}

	DB struct {
		User     string        `yaml:"user"`
		Password string        `env-required:"true" env:"DB_PASSWORD"`
		Name     string        `yaml:"name"`
		Host     string        `yaml:"host"`
		Port     string        `yaml:"port"`
		Timeout  time.Duration `yaml:"timeout"`
		DSN      string        `yaml:"-"`
	}

	Ntfy struct {
		Host  string `yaml:"host"`
		Port  string `yaml:"port"`
		Theme string `yaml:"theme"`
		Token string `env-required:"true" env:"NTFY_TOKEN"`
	}

	App struct {
		LocaleName  string         `yaml:"locale"`
		Locale      *time.Location `yaml:"-"`
		Logging     `yaml:"logging"`
		HealthCheck `yaml:"healthcheck"`
		Climate     `yaml:"climate"`
	}

	Logging struct {
		LogLevel int `yaml:"log_level"`
		// use JSON format if true
		JSONFormat bool `yaml:"json_format"`
	}

	HealthCheck struct {
		Port string `yaml:"port"`
	}

	Climate struct {
		CheckDBEvery time.Duration `yaml:"check_db_every"`
		TempTreshold `yaml:"temp_threshold"`
	}

	TempTreshold struct {
		Standart float32 `yaml:"standart"`
		High     float32 `yaml:"high"`
		Urgent   float32 `yaml:"urgent"`
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
		App: App{
			LocaleName: _defLocaleName,
			Logging: Logging{
				LogLevel:   _defLogLevel,
				JSONFormat: _defJSONFormat,
			},
			HealthCheck: HealthCheck{
				Port: _defHealthcheckPort,
			},
			Climate: Climate{
				CheckDBEvery: _defCheckDBEvery,
				TempTreshold: TempTreshold{
					Standart: _defTempTresholdStandart,
					High:     _defTempTresholdHigh,
					Urgent:   _defTempTresholdUrgent,
				},
			},
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

	// parse locale name
	locale, err := time.LoadLocation(cfg.App.LocaleName)
	if err != nil {
		return nil, fmt.Errorf("parse locale name: %w", err)
	}
	cfg.App.Locale = locale

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

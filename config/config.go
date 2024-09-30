package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App
		JWT
		HTTP
		Log
		DB
	}

	// App -.
	App struct {
		Name    string `env-required:"true" env:"APP_NAME"`
		Version string `env-required:"true" env:"APP_VERSION"`
	}

	JWT struct {
		Secret     string `env-required:"true" env:"JWT_SECRET"`
		Expiration int    `env-required:"true" env:"JWT_EXPIRATION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" env:"PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" env:"LOG_LEVEL"`
	}

	// DB -.
	DB struct {
		URL      string `env:"DB_URL"`
		HOST     string `env:"DB_HOST"`
		PORT     string `env:"DB_PORT"`
		DATABASE string `env:"DB_DATABASE"`
		USERNAME string `env:"DB_USERNAME"`
		PASSWORD string `env:"DB_PASSWORD"`
		TIMEZONE string `env:"DB_TIMEZONE"`
		SSL_MODE string `env:"DB_SSL_MODE"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./.env", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

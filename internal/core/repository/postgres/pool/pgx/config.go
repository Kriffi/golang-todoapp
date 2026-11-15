package core_pgx_pool

import (
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseURL string
	Host        string        `envconfig:"HOST" required:"true"`
	Port        string        `envconfig:"PORT" default:"5432"`
	User        string        `envconfig:"USER" required:"true"`
	Password    string        `envconfig:"PASSWORD" required:"true"`
	Database    string        `envconfig:"DB" required:"true"`
	Timeout     time.Duration `envconfig:"TIMEOUT" required:"true"`
}

func NewConfig() (Config, error) {
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		var options struct {
			Timeout time.Duration `envconfig:"TIMEOUT" default:"10s"`
		}
		if err := envconfig.Process("POSTGRES", &options); err != nil {
			return Config{}, fmt.Errorf("process postgres timeout: %w", err)
		}
		return Config{DatabaseURL: databaseURL, Timeout: options.Timeout}, nil
	}

	var config Config

	if err := envconfig.Process("POSTGRES", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get Postgres connection pool config: %w", err)
		panic(err)
	}

	return config
}

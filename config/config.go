package config

import (
	"os"

	"github.com/go-zoox/ini"
)

const configFile = "config.ini"

type Config struct {
	EjConfig     `ini:"ejudge"`
	DBConfig     `ini:"database"`
	ServerConfig `ini:"server"`
}

type EjConfig struct {
	APIKey    string `ini:"API_KEY"`
	APISecret string `ini:"API_SECRET"`
	URL       string `ini:"URL"`
}

type DBConfig struct {
	SqlitePath string `ini:"SQLITE_PATH"`
}

type ServerConfig struct {
	GinSecret        string `ini:"GIN_SECRET"`
	JWTSecret        string `ini:"JWT_SECRET"`
	PollMaxRuns      int64  `ini:"POLL_MAX_RUNS"`
	PollDelaySeconds int64  `ini:"POLL_DELAY_SECONDS"`
	ReviewLimit      int64  `ini:"REVIEW_LIMIT"`
}

func NewConfig() (*Config, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := ini.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

package config

import (
	"encoding/json"

	"github.com/joho/godotenv"
)

type Config struct {
	EjConfig
	DBConfig
	GinConfig
}

type EjConfig struct {
	APIKey    string `json:"EJ_API_KEY"`
	APISecret string `json:"EJ_API_SECRET"`
	URL       string `json:"EJ_URL"`
}

type DBConfig struct {
	SqlitePath string `json:"SQLITE_PATH"`
}

type GinConfig struct {
	Secret    string `json:"GIN_SECRET"`
	JWTSecret string `json:"JWT_SECRET"`
}

func NewConfig() (*Config, error) {
	env, err := godotenv.Read()
	if err != nil {
		return nil, err
	}
	senv, err := json.Marshal(env)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(senv, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

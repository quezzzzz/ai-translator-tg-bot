package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Postgres     *PostgresConfig
	AITranslator *AITranslatorConfig
	TGBotToken   string
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	Database string `json:"database"`
	SSLMode  string `json:"sslmode"`

	User     string
	Password string

	UserEnvName     string `json:"USER_ENV_NAME"`
	PasswordEnvName string `json:"PASSWORD_ENV_NAME"`
}

type AITranslatorConfig struct {
	URL string `json:"url"`

	APIKey   string
	FolderID string

	APIKeyEnvName   string `json:"api_key_env_name"`
	FolderIDEnvName string `json:"folder_id_env_name"`
}

func ParseConfig(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, config.fillEnvs()
}

func (c *Config) fillEnvs() error {
	c.Postgres.User = os.Getenv(c.Postgres.UserEnvName)
	c.Postgres.Password = os.Getenv(c.Postgres.PasswordEnvName)
	if c.Postgres.User == "" || c.Postgres.Password == "" {
		return fmt.Errorf("bad postgres envs")
	}

	c.AITranslator.APIKey = os.Getenv(c.AITranslator.APIKeyEnvName)
	c.AITranslator.FolderID = os.Getenv(c.AITranslator.FolderIDEnvName)
	if c.AITranslator.FolderID == "" || c.AITranslator.APIKey == "" {
		return fmt.Errorf("bad yandex ai translator envs")
	}
	return nil
}

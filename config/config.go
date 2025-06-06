package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type (
	Config struct {
		DefaultDsn string `json:"defaultDsn"`
	}
)

var userConfigDir, _ = os.UserConfigDir()
var configDir = userConfigDir + "/sharky/"
var configFile = configDir + "config.json"

func SetDefaultDsn(dsn string) error {
	config := Config{dsn}
	fileBytes, _ := json.Marshal(config)
	err := os.WriteFile(configFile, fileBytes, 0600)
	if err != nil {
		return fmt.Errorf("Failed to save default dsn: %w", err)
	}
	return nil
}

func GetDefaultDsn() (string, error) {
	filesBytes, err := os.ReadFile(configFile)
	if err != nil {
		return "", fmt.Errorf("Failed to read config file: %w", err)
	}
	var config Config
	err = json.Unmarshal(filesBytes, &config)
	if err != nil {
		return "", fmt.Errorf("Failed to parse config file: %w", err)
	}
	if len(config.DefaultDsn) == 0 {
		return "", errors.New("No default DSN specified")
	}
	return config.DefaultDsn, nil
}

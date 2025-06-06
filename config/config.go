package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sharky/api"
	"time"
)

type (
	Config struct {
		DefaultDsn string `json:"defaultDsn"`
	}

	authInfo struct {
		AccessToken string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		TimeToLive int `json:"timeToLive"`
		Timestamp string `json:"timestamp"`
	}
)

var userConfigDir, _ = os.UserConfigDir()
var configDir = userConfigDir + "/sharky/"
var authFile = configDir + "auth.json"
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

func SaveAuth(accessToken string, refreshToken string, timeToLive int) error {
	authInfo := authInfo{
		accessToken,
		refreshToken,
		timeToLive,
		time.Now().Format(time.RFC3339),
	}
	fileBytes, _ := json.Marshal(authInfo)
	err := os.WriteFile(authFile, fileBytes, 0600)
	if err != nil {
		return fmt.Errorf("Failed to save auth to file: %w", err)
	}
	return nil
}

func GetAccessToken() (string, error) {
	fileBytes, err := os.ReadFile(authFile)
	if err != nil {
		return "", fmt.Errorf("Failed to read auth file: %w", err)
	}
	var authInfo authInfo
	err = json.Unmarshal(fileBytes, &authInfo)
	if err != nil {
		return "", fmt.Errorf("Failed to parse auth file: %w", err)
	}
	timestamp, _ := time.Parse(time.RFC3339, authInfo.Timestamp)
	expirationTime := timestamp.Add(time.Second * time.Duration(authInfo.TimeToLive))
	if time.Now().After(expirationTime) {
		fmt.Println("Token expired, requesting new one...")
		newTokens, err := api.RefreshToken(authInfo.RefreshToken)
		if err != nil {
			return "", err
		}
		SaveAuth(newTokens.AccessToken, newTokens.RefreshToken, newTokens.TimeToLive)
	}
	return authInfo.AccessToken, nil
}

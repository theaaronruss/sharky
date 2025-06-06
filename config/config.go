package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sharky/api"
	"time"
)

const (
	DefaultDsn = "defaultDsn"
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

func SetProperty(name string, value string) error {
	fileBytes, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("Failed to read config file: %w", err)
	}
	var config map[string]string;
	err = json.Unmarshal(fileBytes, &config)
	if err != nil {
		return fmt.Errorf("Failed to parse config file: %w", err)
	}
	config[name] = value
	fileBytes, _ = json.Marshal(config)
	err = os.WriteFile(configFile, fileBytes, 0600)
	if err != nil {
		return fmt.Errorf("Failed to save config file: %w", err)
	}
	return nil
}

func GetProperty(name string) (string, error) {
	fileBytes, err := os.ReadFile(configFile)
	if err != nil {
		return "", fmt.Errorf("Failed to read config file: %w", err)
	}
	var config map[string]string
	err = json.Unmarshal(fileBytes, &config)
	if err != nil {
		return "", fmt.Errorf("Failed to parse config file: %w", err)
	}
	value := config[name]
	if value == "" {
		return value, fmt.Errorf("Config property %s not found", name)
	}
	return config[name], nil
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

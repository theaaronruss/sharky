package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type (
	Tokens struct {
		AccessToken string
		RefreshToken string
		TimeToLive int
	}

	refreshTokenResponse struct {
		AccessToken string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn int `json:"expires_in"`
	}
)

func GenerateLoginUrl() (string, string) {
	url := "https://login.sharkninja.com/authorize"
	url += "?ui_locales=en"
	url += "&response_type=code"
	url += "&redirect_uri=com.sharkninja.shark://login.sharkninja.com/ios/com.sharkninja.shark/callback"
	url += "&screen_hint=signin"
	url += "&code_challenge_method=S256"
	hasher := sha256.New()
	codeVerifier := generateRandomString(43)
	hasher.Write([]byte(codeVerifier))
	hashBytes := hasher.Sum(nil)
	codeChallenge := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hashBytes)
	url += "&code_challenge=" + codeChallenge
	url += "&client_id=wsguxrqm77mq4LtrTrwg8ZJUxmSrexGi"
	url += "&state=" + generateRandomString(43)
	url += "&scope=openid%20profile%20email%20offline_access%20read:users%20read:current_user%20read:user_idp_tokens"
	url += "&auth0Client=eyJlbnYiOnsiaU9TIjoiMTguNSIsInN3aWZ0IjoiNS54In0sIm5hbWUiOiJBdXRoMC5zd2lmdCIsInZlcnNpb24iOiIyLjYuMCJ9"
	return url, codeVerifier
}

func generateRandomString(length int) string {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	output := make([]byte, length)
	for i := range length {
		randIdx := rand.Intn(len(chars))
		output[i] = chars[randIdx]
		length--
	}
	return string(output)
}

func RefreshToken(refreshToken string) (Tokens, error) {
	requestBody := []byte("{\"user\":{\"refresh_token\": \"" + refreshToken + "\"}}")
	request, err := http.NewRequest(http.MethodPost, "https://user-field-39a9391a.aylanetworks.com/users/refresh_token.json", bytes.NewReader(requestBody))
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		return Tokens{}, fmt.Errorf("Failed to create request for refreshing token: %w", err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return Tokens{}, fmt.Errorf("Failed to refresh token: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode == 401 {
		return Tokens{}, errors.New("Invalid refresh token")
	}
	var newTokens refreshTokenResponse
	responseBytes, err := io.ReadAll(response.Body)
	if json.Unmarshal(responseBytes, &newTokens) != nil {
		return Tokens{}, fmt.Errorf("Failed to parse refresh token response: %w", err)
	}
	time.Sleep(time.Millisecond * 250)
	return Tokens{
		AccessToken: newTokens.AccessToken,
		RefreshToken: newTokens.RefreshToken,
		TimeToLive: newTokens.ExpiresIn,
	}, nil
}

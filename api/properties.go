package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type (
	Property struct {
		Name string
		Dsn string
		Value string
	}

	getDatapointResponse struct {
		Datapoint struct {
			Value json.RawMessage `json:"value"`
		} `json:"datapoint"`
	}
)

func BatchUpdateProperties(accessToken string, userUuid string, properties ...Property) error {
	if len(properties) < 1 {
		return nil
	}
	requestBody := "{\"batch_datapoints\":["
	for i, property := range properties {
		requestBody += "{\"datapoint\":{\"metadata\":{\"userUUID\":\"" + userUuid + "\"},\"value\":\"" + property.Value + "\"},\"dsn\":\"" + property.Dsn + "\",\"name\":\"" + property.Name + "\"}"
		if i != len(properties) - 1 {
			requestBody += ","
		}
	}
	requestBody += "]}"
	request, err := http.NewRequest(http.MethodPost, "https://ads-field-39a9391a.aylanetworks.com/apiv1/batch_datapoints.json", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return fmt.Errorf("Failed to create request to set batch properties: %w", err)
	}
	request.Header.Add("Authorization", accessToken)
	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("Failed to set batch properties: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode == 401 {
		return errors.New("Invalid response while adding batch properties: not authorized")
	} else if response.StatusCode != 201 {
		return errors.New("Invalid response while adding batch properties: http status " + fmt.Sprint(response.StatusCode))
	}
	return nil
}

func UpdateProperty(accessToken string, userUuid string, dsn string, name string, value string) error {
	requestBody := "{\"datapoint\":{\"metadata\":{\"userUUID\":\"" + userUuid + "\"},\"value\":\"" + value + "\"}}"
	request, err := http.NewRequest(http.MethodPost, "https://ads-field-39a9391a.aylanetworks.com/apiv1/dsns/" + dsn + "/properties/" + name + "/datapoints.json", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return fmt.Errorf("Failed to create request to set property: %w", err)
	}
	request.Header.Add("Authorization", accessToken)
	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("Failed to set property: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode == 401 {
		return errors.New("Invalid response while adding property: not authorized")
	} else if response.StatusCode != 201 {
		return errors.New("Invalid response while adding property: http status: " + fmt.Sprint(response.StatusCode))
	}
	return nil
}

func GetLatestDatapointValue(accessToken string, dsn string, name string) (string, error) {
	request, err := http.NewRequest(http.MethodGet, "https://ads-field-39a9391a.aylanetworks.com/apiv1/dsns/" + dsn + "/properties/" + name + "/datapoints.json?limit=1", nil)
	if err != nil {
		return "", fmt.Errorf("Failed to create request to get datapoint: %w", err)
	}
	request.Header.Add("Authorization", accessToken)
	request.Header.Add("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("Failed to get datapoint: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode == 401 {
		return "", errors.New("Invalid response while getting datapoint: not authorized")
	} else if response.StatusCode == 404 {
		return "", errors.New("Invalid response while getting datapoint: invalid DSN or property name")
	}
	var datapoints []getDatapointResponse
	responseBytes, _ := io.ReadAll(response.Body)
	json.Unmarshal(responseBytes, &datapoints)
	if len(datapoints) == 0 {
		return "", errors.New("No datapoints found for property")
	}
	value := string(datapoints[0].Datapoint.Value)
	value = strings.ReplaceAll(value, "\"", "")
	return string(value), nil
}

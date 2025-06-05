package api

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

type (
	Property struct {
		Name string
		Dsn string
		Value string
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

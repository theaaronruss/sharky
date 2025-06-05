package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type (
	Device struct {
		Name string
		Status string
	}

	DeviceListResponse struct {
		Device struct {
			ProductName string `json:"product_name"`
			ConnectionStatus string `json:"connection_status"`
		} `json:"device"`
	}
)

func GetDeviceList(accessToken string) ([]Device, error) {
	request, err := http.NewRequest(http.MethodGet, "https://ads-field-39a9391a.aylanetworks.com/apiv1/devices.json", nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request for retrieving device list: %w", err)
	}
	request.Header.Add("Authorization", accessToken)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve device list: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode == 401 {
		return nil, errors.New("Invalid response while retrieving device list: not authorized")
	} else if response.StatusCode == 403 {
		return nil, errors.New("Invalid response while retrieving device list: forbidden")
	}
	var responseBody []DeviceListResponse
	responseBytes, err := io.ReadAll(response.Body)
	if json.Unmarshal(responseBytes, &responseBody) != nil {
		return nil, fmt.Errorf("Failed to parse device list response: %w", err)
	}
	deviceList := make([]Device, len(responseBody))
	for i, device := range responseBody {
		deviceList[i] = Device{
			device.Device.ProductName,
			device.Device.ConnectionStatus,
		}
	}
	return deviceList, nil
}

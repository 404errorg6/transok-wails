package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"transok/backend/consts"
	"transok/backend/services"
)

type DiscoverDevice struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Uname    string `json:"uname"`
	Platform string `json:"platform"`
	Address  string `json:"address"`
}

type DiscoverHandler struct {
}

var (
	discoverHandler *DiscoverHandler
	once            sync.Once
)

func GetDiscoverHandler() *DiscoverHandler {
	once.Do(func() {
		discoverHandler = &DiscoverHandler{}
	})
	return discoverHandler
}

func (h *DiscoverHandler) GetType() string {
	return "DISCOVER"
}

func (h *DiscoverHandler) Handle(payload consts.DiscoverPayload) {
	device := DiscoverDevice{
		IP:       payload.Payload["IP"],
		Port:     payload.Payload["Port"],
		Uname:    payload.Payload["Uname"],
		Platform: payload.Payload["Platform"],
		Address:  fmt.Sprintf("%s:%s", payload.Payload["IP"], payload.Payload["Port"]),
	}

	var deviceList []DiscoverDevice
	if discoverList, ok := services.Storage().Get("discover-list"); ok {
		jsonData, _ := json.Marshal(discoverList)
		json.Unmarshal(jsonData, &deviceList)
	}

	// Deduplicate by Address; update information if it already exists
	found := false
	for i, dev := range deviceList {
		if dev.Address == device.Address {
			deviceList[i] = device // Update device info
			found = true
			break
		}
	}
	if !found {
		deviceList = append(deviceList, device)
	}

	services.Storage().Set("discover-list", deviceList)
}

/* Test if the address is reachable */
func ping(address string) bool {
	resp, err := http.Get("http://" + address + "/discover/ping")
	if err != nil {
		fmt.Printf("ping %s failed: %v\n", address, err)
		return false
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read response: %v\n", err)
		return false
	}

	// Parse JSON response
	var result struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("failed to parse JSON: %v\n", err)
		return false
	}

	return result.Success
}

func mapperDiscoverDevice(device map[string]interface{}) DiscoverDevice {
	return DiscoverDevice{
		IP:       device["ip"].(string),
		Port:     device["port"].(string),
		Uname:    device["uname"].(string),
		Platform: device["platform"].(string),
		Address:  device["address"].(string),
	}
}

/* Get discovery list */
func (h *DiscoverHandler) GetDiscoverList() []DiscoverDevice {
	var deviceList []DiscoverDevice
	if discoverList, ok := services.Storage().Get("discover-list"); ok {
		jsonData, _ := json.Marshal(discoverList)
		json.Unmarshal(jsonData, &deviceList)
	}

	// Return only reachable devices and remove inaccessible ones from storage
	var accessibleDevices []DiscoverDevice
	var needUpdate bool

	for _, dev := range deviceList {
		if ping(dev.Address) {
			accessibleDevices = append(accessibleDevices, dev)
		} else {
			needUpdate = true // Mark storage for update
		}
	}

	// Update storage if any inaccessible devices were filtered out
	if needUpdate {
		services.Storage().Set("discover-list", accessibleDevices)
	}

	return accessibleDevices
}

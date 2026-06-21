package guard

import (
	"encoding/json"
	"os"
)

type AppConfig struct {
	AllowedNetInterfaces []string `json:"allowed_net_interfaces"`
	AllowedBluetoothMACs []string `json:"allowed_bluetooth_macs"`
	AllowedUSBSerials    []string `json:"allowed_usb_serials"`
}

func LoadAppConfig() (*AppConfig, error) {
	file, err := os.ReadFile("../config.json")
	if err != nil {
		return nil, err
	}

	var config AppConfig
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

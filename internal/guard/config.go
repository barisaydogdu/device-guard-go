package guard

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
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

func LoadAppConfigTxt() (*AppConfig, error) {
	config := &AppConfig{}

	file, err := os.Open("../guard.conf")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			continue
		}

		values := strings.Split(line, "=")

		if len(values) != 2 {
			log.Println("values different 2")
			continue
		}

		switch values[0] {
		case "net":
			config.AllowedNetInterfaces = append(config.AllowedNetInterfaces, values[1])
		case "bt":
			config.AllowedBluetoothMACs = append(config.AllowedBluetoothMACs, values[1])
		case "usb":
			config.AllowedUSBSerials = append(config.AllowedUSBSerials, values[1])
		}
	}
	return config, scanner.Err()
}

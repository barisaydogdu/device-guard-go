package guard

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type AppConfig struct {
	mu                   sync.Mutex
	ConfigFilePath       string
	AllowedNetInterfaces []string `json:"allowed_net_interfaces"`
	AllowedBluetoothMACs []string `json:"allowed_bluetooth_macs"`
	AllowedUSBSerials    []string `json:"allowed_usb_serials"`
	AllowedCameraIDsb    []string `json:"allowed_camera_ids"`
}

const guardConfigFile = "../guard.conf"

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

func LoadAppConfigTxt(configPath string) (*AppConfig, error) {
	config := &AppConfig{}

	file, err := os.Open(filepath.Join("..", configPath))
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	scanner := bufio.NewScanner(file)

	config.ConfigFilePath = configPath

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
		break
	}
	return config, scanner.Err()
}

func (c *AppConfig) AddAllowedDevice(deviceType string, identifier string) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	var logPrefix string
	var targetHeader string

	cleanedIdentifier := strings.TrimSpace(identifier)

	switch deviceType {
	case "usb":
		c.AllowedUSBSerials = append(c.AllowedUSBSerials, cleanedIdentifier)
		logPrefix = "usb"
		targetHeader = "# === USB Guard Permissions ==="
	case "blueetooth":
		c.AllowedBluetoothMACs = append(c.AllowedBluetoothMACs, cleanedIdentifier)
		logPrefix = "bt"
		targetHeader = "# === Bluetooth Guard Permissions ==="
	case "net":
		c.AllowedNetInterfaces = append(c.AllowedNetInterfaces, cleanedIdentifier)
		logPrefix = "net"
		targetHeader = "# === Network Guard Permissions ==="
	default:
		return fmt.Errorf("[config] unknown device type %s", deviceType)
	}

	content, err := os.ReadFile(guardConfigFile)
	if err != nil {
		return fmt.Errorf("[config] config file cannot read")
	}

	lines := strings.Split(string(content), "\n")

	tIndex := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == targetHeader {
			tIndex = i
			break
		}
	}

	if tIndex == -1 {
		return fmt.Errorf("[config] title '%s' not found in file", targetHeader)
	}

	insertIndex := tIndex + 1
	lineToWrite := fmt.Sprintf("%s=%s", logPrefix, cleanedIdentifier)

	lines = append(lines[:insertIndex], append([]string{lineToWrite}, lines[insertIndex:]...)...)

	finalContent := strings.Join(lines, "\n")
	err = os.WriteFile(guardConfigFile, []byte(finalContent), 0644)
	if err != nil {
		return fmt.Errorf("[config] config Failed to update configuration file:  %w", err)
	}

	return nil
}

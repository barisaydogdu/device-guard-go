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

	defer file.Close()

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
			break
		case "bt":
			config.AllowedBluetoothMACs = append(config.AllowedBluetoothMACs, values[1])
			break
		case "usb":
			config.AllowedUSBSerials = append(config.AllowedUSBSerials, values[1])
			break
		}
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
		fmt.Println("qqq1")
		break
	case "blueetooth":
		c.AllowedBluetoothMACs = append(c.AllowedBluetoothMACs, cleanedIdentifier)
		logPrefix = "bt"
		targetHeader = "# === Bluetooth Guard Permissions ==="
		break
	case "net":
		c.AllowedNetInterfaces = append(c.AllowedNetInterfaces, cleanedIdentifier)
		logPrefix = "net"
		targetHeader = "# === Network Guard Permissions ==="
		break
	default:
		return fmt.Errorf("unknown device type %s", deviceType)
	}

	content, err := os.ReadFile("../guard.conf")
	if err != nil {
		return fmt.Errorf("config file cannot read")
	}

	lines := strings.Split(string(content), "\n")

	var tIndex int = -1
	for i, line := range lines {
		if strings.TrimSpace(line) == targetHeader {
			tIndex = i
			break
		}
	}

	if tIndex == -1 {
		return fmt.Errorf("kritik hata: '%s' basligi dosyada bulunamadi", targetHeader)
	}

	insertIndex := tIndex + 1
	lineToWrite := fmt.Sprintf("%s=%s", logPrefix, cleanedIdentifier)

	lines = append(lines[:insertIndex], append([]string{lineToWrite}, lines[insertIndex:]...)...)

	finalContent := strings.Join(lines, "\n")
	err = os.WriteFile("../guard.conf", []byte(finalContent), 0644)
	if err != nil {
		return fmt.Errorf("config dosyasi guncellenemedi: %w", err)
	}

	return nil
}

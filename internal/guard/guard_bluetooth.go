package guard

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type BluetoothGuard struct {
	Config *AppConfig
}

func (b *BluetoothGuard) Block(eventMap map[string]string) error {
	if eventMap["DEVTYPE"] == "host" {
		return nil
	}

	for key, value := range eventMap {
		log.Printf("%s = %s\n", key, value)
	}

	addressPath := filepath.Join("/sys", eventMap["DEVPATH"], "address")
	address, err := os.ReadFile(addressPath)
	if err != nil {
		return err
	}

	macAddress := strings.TrimSpace(string(address))

	log.Println("mac address:", macAddress)

	cmd := exec.Command("rfkill", "block", "bluetooth")
	err = cmd.Run()
	if err != nil {
		log.Println("Failed to block device:", err)
		return err
	}

	log.Println("blocked blueetooth")
	return nil
}

func (b *BluetoothGuard) Allow(devPath string) error {
	cmd := exec.Command("rfkill", "unblock", "bluetooth")
	err := cmd.Run()
	if err != nil {
		log.Println("Failed to unblock device:", err)
		return err
	}
	return nil
}

func (b *BluetoothGuard) BlockSpecificDevice(macAddress string) error {
	if macAddress == "" {
		return fmt.Errorf("bluetooth mac address is empty")
	}

	exec.Command("bluetoothctl", "disconnect", macAddress).Run()

	cmd := exec.Command("bluetoothctl", "block", macAddress)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cihaz cannot block: %v, detail: %s", err, string(output))
	}

	log.Println("blocked blueetooth successfully")
	return nil
}

func (b *BluetoothGuard) AllowSpecificDevice(macAddress string) error {
	if macAddress == "" {
		return fmt.Errorf("bluetooth mac address is empty")
	}

	cmd := exec.Command("bluetoothctl", "unblock", macAddress)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cihazın engeli kaldırılamadı: %v, detay: %s", err, string(output))
	}

	log.Println("unblocked blueetooth successfully")
	return nil
}

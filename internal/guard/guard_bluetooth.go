package guard

import (
	"log"
	"os/exec"
)

type BluetoothGuard struct{}

func (b *BluetoothGuard) Block(eventMap map[string]string) error {
	if eventMap["DEVTYPE"] == "host" {
		return nil
	}

	cmd := exec.Command("rfkill", "block", "bluetooth")
	err := cmd.Run()
	if err != nil {
		log.Println("Failed to block device:", err)
		return err
	}

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

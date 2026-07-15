package guard

import (
	"fmt"
	"github/usb-guard-go/internal/util"
	"log"
	"os/exec"
	"strings"
)

type BluetoothGuard struct {
	Config *AppConfig
}

func (b *BluetoothGuard) Block(eventMap map[string]string) error {
	// if eventMap["DEVTYPE"] == "host" {
	// 	return nil
	// }

	// for key, value := range eventMap {
	// 	log.Printf("%s = %s\n", key, value)
	// }

	// addressPath := filepath.Join("/sys", eventMap["DEVPATH"], "address")
	// address, err := os.ReadFile(addressPath)
	// if err != nil {
	// 	return err
	// }

	// macAddress := strings.TrimSpace(string(address))

	// for i := 0; i < len(b.Config.AllowedBluetoothMACs); i++ {
	// 	if b.Config.AllowedBluetoothMACs[i] == macAddress {
	// 		log.Printf("[Blueetooth] mac address is whitelisted. Allowing connection.\n")
	// 		return nil
	// 	}
	// }

	// cmd := exec.Command("rfkill", "block", "bluetooth")
	// err = cmd.Run()
	// if err != nil {
	// 	log.Println("Failed to block device:", err)
	// 	return err
	// }

	// log.Println("blocked blueetooth")

	// util.SendNotification("Device Blocked", fmt.Sprintf("An unauthorized device was detected and blocked.\nDevice ID: %s", macAddress))

	// return nil
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
	fmt.Println("ccc2")
	if macAddress == "" {
		return fmt.Errorf("bluetooth mac address is empty")
	}

	for i := 0; i < len(b.Config.AllowedBluetoothMACs); i++ {
		if b.Config.AllowedBluetoothMACs[i] == strings.TrimSpace(string(macAddress)) {
			log.Printf("[USB] serial Number is whitelisted. Allowing connection.\n")
			return nil
		}
	}

	exec.Command("bluetoothctl", "disconnect", macAddress).Run()

	cmd := exec.Command("bluetoothctl", "block", macAddress)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("device cannot block: %v, detail: %s", err, string(output))
	}

	log.Println("blocked blueetooth successfully")

	util.SendNotification("Device Blocked", fmt.Sprintf("An unauthorized device was detected and blocked.\nDevice ID: %s", macAddress))

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

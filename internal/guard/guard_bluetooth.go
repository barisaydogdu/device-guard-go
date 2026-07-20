package guard

import (
	"fmt"
	"github/usb-guard-go/internal/util"
	"log"
	"os/exec"
	"strings"
	"sync"
)

type BluetoothGuard struct {
	Config  *AppConfig
	pending map[string]bool
	mu      sync.Mutex
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
	if macAddress == "" {
		return fmt.Errorf("bluetooth mac address is empty")
	}

	b.mu.Lock()

	if b.pending == nil {
		b.pending = make(map[string]bool)
	}

	if b.pending[macAddress] {
		b.mu.Unlock()
		return nil
	}

	b.pending[macAddress] = true
	b.mu.Unlock()

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

	go func(macAddress string) {
		defer func() {
			b.mu.Lock()
			delete(b.pending, macAddress)
			b.mu.Unlock()
		}()

		if util.AskPermission("Device Blocked", fmt.Sprintf("An unauthorized device was detected and blocked.\nDevice ID: %s", macAddress)) {
			err := b.Config.AddAllowedDevice("blueetooth", macAddress)
			if err != nil {
				log.Printf("[bluetooth] Failed to add allowed blueetooth device: %v", err)
				return
			}
			err = b.AllowSpecificDevice(macAddress)
			if err != nil {
				log.Printf("[blueetooth] Failed to allow device: %v", err)
				return
			}
		}
	}(macAddress)

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

	exec.Command("bluetoothctl", "trust", macAddress).Run()

	cmdConnect := exec.Command("bluetoothctl", "connect", macAddress)
	_, connErr := cmdConnect.CombinedOutput()
	if connErr != nil {
		log.Printf("[Bluetooth] Device added to the list, but failed to connect automatically.")
	} else {
		log.Printf("[Bluetooth] Automatic connection to the device was successful.")
	}

	log.Println("unblocked and trusted blueetooth successfully")
	return nil
}

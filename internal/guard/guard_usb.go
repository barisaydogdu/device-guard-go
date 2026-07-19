package guard

import (
	"fmt"
	"github/usb-guard-go/internal/util"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type USBGuard struct {
	Config *AppConfig
}

func (u *USBGuard) Block(eventMap map[string]string) error {
	if eventMap["DEVTYPE"] != "usb_device" {
		return nil
	}

	serialNumberPath := filepath.Join("/sys"+eventMap["DEVPATH"], "/serial")
	serial, err := os.ReadFile(serialNumberPath)
	if err != nil {
		log.Println("serial number not found suspicious device")
	}

	for i := 0; i < len(u.Config.AllowedUSBSerials); i++ {
		if u.Config.AllowedUSBSerials[i] == strings.TrimSpace(string(serial)) {
			log.Printf("[USB] serial Number is whitelisted. Allowing connection.\n")
			return nil
		}
	}

	targetFile := filepath.Join("/sys", eventMap["DEVPATH"], "authorized")

	err = os.WriteFile(targetFile, []byte("0"), 0644)
	if err != nil {
		return err
	}

	log.Println("Device blocked successfully")

	go func(blockSerial string, authFile string) {
		if util.AskPermission("Device Blocked", fmt.Sprintf("An unauthorized device was detected and blocked.\nDevice ID: %s", serial)) {
			log.Printf("Kullanici %s seri numarali cihaza izin verdi!", blockSerial)
			u.Config.AddAllowedDevice("usb", string(serial))

			err := os.WriteFile(targetFile, []byte("1"), 0644)
			if err != nil {
				log.Printf("[ERROR] Device added to the list, but the kernel lock could not be unlocked: %v\n", err)
			} else {
				log.Println("The device was woken up and connected to the system without physically unplugging and replugging it!")
			}
		} else {
			log.Printf("The user ignored or closed the notification.")
		}

	}(string(serial), targetFile)
	return nil
}

func (u *USBGuard) Allow(devPath string) error {
	targetFile := filepath.Join("/sys", devPath, "authorized")

	err := os.WriteFile(targetFile, []byte("1"), 0644)
	if err != nil {
		return err
	}

	log.Println("Device allowed successfully")

	return nil
}

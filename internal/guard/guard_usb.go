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

	util.SendNotification("Device Blocked", fmt.Sprintf("An unauthorized device was detected and blocked.\nDevice ID: %s", serial))

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

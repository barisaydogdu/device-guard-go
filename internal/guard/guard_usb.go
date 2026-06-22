package guard

import (
	"fmt"
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
	file, err := os.ReadFile(serialNumberPath)
	if err != nil {
		log.Println("serial number not found suspicious device")
	}

	fmt.Print(string(file))
	for i := 0; i < len(u.Config.AllowedUSBSerials); i++ {
		if u.Config.AllowedUSBSerials[i] == strings.TrimSpace(string(file)) {
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

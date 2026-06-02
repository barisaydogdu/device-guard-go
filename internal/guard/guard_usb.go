package guard

import (
	"log"
	"os"
	"path/filepath"
)

type USBGuard struct{}

func (u *USBGuard) Block(eventMap map[string]string) error {
	if eventMap["DEVTYPE"] != "usb_device" {
		return nil
	}

	log.Println("USB Device blocking")
	log.Println("Product ID:", eventMap["PRODUCT"])

	targetFile := filepath.Join("/sys", eventMap["DEVPATH"], "authorized")

	err := os.WriteFile(targetFile, []byte("0"), 0644)
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

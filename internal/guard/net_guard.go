package guard

import (
	"log"
	"os/exec"
	"path/filepath"
)

type NetGuard struct{}

func (b *NetGuard) Block(eventMap map[string]string) error {
	iFace := eventMap["INTERFACE"]

	if iFace != "" {
		err := exec.Command("ip", "link", "set", "dev", iFace, "down").Run()
		if err != nil {
			log.Println("Failed to block net:", err)
			return err
		}
		log.Println("Successfully blocked net:", iFace)
	}

	return nil
}

func (b *NetGuard) Allow(devPath string) error {
	ipFace := filepath.Base(devPath)
	if ipFace != "" {
		err := exec.Command("ip", "link", "set", "dev", ipFace, "up").Run()
		if err != nil {
			log.Println("Failed to allow net:", err)
			return err
		}
	}
	return nil
}

package guard

import (
	"fmt"
	"github/usb-guard-go/internal/util"
	"log"
	"os/exec"
	"path/filepath"
)

type NetGuard struct {
	Config *AppConfig
}

func (b *NetGuard) Block(eventMap map[string]string) error {
	iFace := eventMap["INTERFACE"]

	log.Println("allowed net interface:", b.Config.AllowedNetInterfaces)

	for i := 0; i < len(b.Config.AllowedNetInterfaces); i++ {
		if b.Config.AllowedNetInterfaces[i] == iFace {
			log.Printf("[Net] iFace is whitelisted. Allowing connection.\n")
			return nil
		}
	}
	if iFace != "" {
		err := exec.Command("ip", "link", "set", "dev", iFace, "down").Run()
		if err != nil {
			log.Println("Failed to block net:", err)
			return err
		}
		log.Println("Successfully blocked net:", iFace)

		_ = util.AskPermission("Device Blocked", fmt.Sprintf("An unauthorized device was detected and blocked.\nDevice ID: %s", iFace))
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

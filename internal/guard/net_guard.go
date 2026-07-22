package guard

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/barisaydogdu/device-guard-go/internal/util"
)

type NetGuard struct {
	Config *AppConfig
}

func (b *NetGuard) Block(eventMap map[string]string) error {
	iFace := eventMap["INTERFACE"]

	log.Println("[net] allowed net interface:", b.Config.AllowedNetInterfaces)

	for i := 0; i < len(b.Config.AllowedNetInterfaces); i++ {
		if b.Config.AllowedNetInterfaces[i] == iFace {
			log.Printf("[net] iFace is whitelisted. Allowing connection.\n")
			return nil
		}
	}
	if iFace != "" {
		err := exec.Command("ip", "link", "set", "dev", iFace, "down").Run()
		if err != nil {
			log.Println("[net] failed to block net:", err)
			return err
		}
		log.Println("[net] successfully blocked net:", iFace)

		go func(iFace string) {
			if util.AskPermission("Net Blocked", fmt.Sprintf("An unauthorized net was detected and blocked.\nDevice ID: %s", iFace)) {
				log.Printf("[net] user allowed net in id: %s", string(iFace))
				err := b.Config.AddAllowedDevice("net", string(iFace))
				if err != nil {
					return
				}
				err = exec.Command("ip", "link", "set", "dev", iFace, "up").Run()
				if err != nil {
					log.Println("[net] failed to allow net:", err)
				}
			}
		}(iFace)
	}

	return nil
}

func (b *NetGuard) Allow(devPath string) error {
	ipFace := filepath.Base(devPath)
	if ipFace != "" {
		err := exec.Command("ip", "link", "set", "dev", ipFace, "up").Run()
		if err != nil {
			log.Println("[net] failed to allow net:", err)
			return err
		}
	}
	return nil
}

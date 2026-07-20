package guard

import (
	"fmt"
	"log"
)

type DeviceBlocker interface {
	Block(eventMap map[string]string) error
	Allow(devPath string) error
}

var guardRegistry = make(map[string]DeviceBlocker)

func RegisterGuard(subSystem string, deviceBlocker DeviceBlocker) {
	guardRegistry[subSystem] = deviceBlocker
}

func HandleDeviceEvent(eventMap map[string]string) error {
	subSystem := eventMap["SUBSYSTEM"]

	if eventMap["ACTION"] != "add" {
		return nil
	}

	if subSystem == "" {
		log.Println("[guard] no subsystem found")
		return nil
	}

	fmt.Printf("Subsystem: %s\n", subSystem)

	guard, exist := guardRegistry[subSystem]

	if !exist {
		log.Println("[guard] no guard found")
		return nil
	}

	log.Printf("guard: %s\n", subSystem)
	return guard.Block(eventMap)
	//return guard.Allow(eventMap["DEVPATH"])
}

func HandleMacEvent(macAddress string) error {
	if macAddress == "" {
		_ = fmt.Errorf("[guard] mac address is empty")
	}

	guard, exist := guardRegistry["bluetooth"]
	if !exist {
		return fmt.Errorf("[guard] no guard found")
	}

	bg, ok := guard.(*BluetoothGuard)
	if !ok {
		return fmt.Errorf("[guard] not a Bluetooth guard")
	}

	log.Printf("guard: %s\n", guard)

	return bg.BlockSpecificDevice(macAddress)
}

package guard

import "log"

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
		log.Println("No subsystem found")
		return nil
	}

	guard, exist := guardRegistry[subSystem]

	if !exist {
		log.Println("No guard found")
		return nil
	}

	log.Printf("Guard: %s\n", subSystem)
	return guard.Block(eventMap)
	//return guard.Allow(eventMap["DEVPATH"])
}

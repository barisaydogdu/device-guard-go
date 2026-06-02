package main

import (
	"github/usb-guard-go/internal/guard"
	kk "github/usb-guard-go/internal/listener"
	"log"
)

func main() {
	guard.RegisterGuard("usb", &guard.USBGuard{})
	guard.RegisterGuard("bluetooth", &guard.BluetoothGuard{})

	_, err := kk.ListenUEvents()
	if err != nil {
		log.Println(err)
		return
	}
}

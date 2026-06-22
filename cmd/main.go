package main

import (
	"context"
	"github/usb-guard-go/internal/guard"
	kk "github/usb-guard-go/internal/listener"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var err error

	config, err := guard.LoadAppConfigTxt()
	if err != nil {
		log.Println("config bulunamadi", err)
		return
	}

	guard.RegisterGuard("usb", &guard.USBGuard{Config: config})
	guard.RegisterGuard("bluetooth", &guard.BluetoothGuard{Config: config})
	guard.RegisterGuard("net", &guard.NetGuard{Config: config})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Println("Starting HCI Listener")

	go func() {
		_, err = kk.HCIListener(ctx)
		if err != nil {
			log.Println(err)
			return
		}

	}()

	go func() {
		_, err = kk.ListenUEvents(ctx)
		if err != nil {
			log.Println(err)
			return
		}
	}()

	log.Println("Waiting for HCI Listener to finish")
	<-ctx.Done()
	log.Println("shutting down")
}

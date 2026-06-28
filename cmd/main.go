package main

import (
	"context"
	"github/usb-guard-go/internal/guard"
	listener "github/usb-guard-go/internal/listener"
	"github/usb-guard-go/internal/util"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var err error

	configPath := util.ResolveConfigPath()

	config, err := guard.LoadAppConfigTxt(configPath)
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
		_, err = listener.HCIListener(ctx)
		if err != nil {
			log.Println(err)
			return
		}

	}()

	go func() {
		_, err = listener.ListenUEvents(ctx)
		if err != nil {
			log.Println(err)
			return
		}
	}()

	log.Println("Waiting for HCI Listener to finish")
	<-ctx.Done()
	log.Println("shutting down")
}

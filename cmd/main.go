package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/barisaydogdu/device-guard-go/internal/guard"
	listener "github.com/barisaydogdu/device-guard-go/internal/listener"
	"github.com/barisaydogdu/device-guard-go/internal/util"
)

func main() {
	var err error

	configPath := util.ResolveConfigPath()

	config, err := guard.LoadAppConfigTxt(configPath)
	if err != nil {
		log.Println("config cannot find", err)
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

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
	guard.RegisterGuard("usb", &guard.USBGuard{})
	guard.RegisterGuard("bluetooth", &guard.BluetoothGuard{})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var err error
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

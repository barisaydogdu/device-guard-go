package main

import (
	kk "github/usb-guard-go/internal"
	"log"
)

func main() {
	_, err := kk.ListenUEvents()
	if err != nil {
		log.Println(err)
		return
	}
}

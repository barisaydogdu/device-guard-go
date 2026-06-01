package internal

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func isUserRoot() bool {
	if os.Getuid() != 0 {
		log.Println("isUserRoot is false")
		return false
	}
	return true
}

func openUEventSocket() (int, error) {
	fd, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW, syscall.NETLINK_KOBJECT_UEVENT)
	if err != nil {
		return 0, err
	}

	socketAddr := syscall.SockaddrNetlink{
		Family: syscall.AF_NETLINK,
		Pad:    0,
		Pid:    0,
		Groups: 1, //uevent listen
	}

	if err := syscall.Bind(fd, &socketAddr); err != nil {
		return 0, err
	}
	return fd, nil
}

func ListenUEvents() (map[string]string, error) {
	fd, err := openUEventSocket()
	if err != nil {
		return nil, err
	}

	defer syscall.Close(fd)

	buffer := make([]byte, 4096)

	for {
		read, err := syscall.Read(fd, buffer)
		if err != nil {
			return nil, fmt.Errorf("socket read error: %v", err)
		}

		parts := bytes.Split(buffer[:read], []byte{0})

		eventMap := make(map[string]string)

		for _, part := range parts {
			if len(part) == 0 {
				continue
			}

			kk := strings.SplitN(string(part), "=", 2)

			if len(kk) == 2 {
				eventMap[kk[0]] = kk[1]
			}
		}

		if eventMap["ACTION"] == "add" && eventMap["SUBSYSTEM"] == "usb" && eventMap["DEVTYPE"] == "usb_device" {
			log.Println("\n--- catch new usb signal---")
			log.Printf("Action:    %s\n", eventMap["ACTION"])
			log.Printf("Sysfs Path: %s\n", eventMap["DEVPATH"])
			log.Printf("ProductID   : %s\n", eventMap["PRODUCT"])

			log.Println("catch usb file")
			err := os.WriteFile(filepath.Join("/sys", eventMap["DEVPATH"], "authorized"), []byte("0"), 0600) // 0644?
			if err != nil {
				return nil, err
			}
		}
	}
}

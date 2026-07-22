package listener

import (
	"context"
	"fmt"
	"log"

	"github.com/barisaydogdu/device-guard-go/internal/guard"

	"golang.org/x/sys/unix"
)

func HCIListener(ctx context.Context) (string, error) {
	log.Println("Starting HCI Listener")
	fd, err := unix.Socket(unix.AF_BLUETOOTH, unix.SOCK_RAW, unix.BTPROTO_HCI)
	if err != nil {
		return "s", err
	}

	defer func(fd int) {
		err := unix.Close(fd)
		if err != nil {
			return
		}
	}(fd)

	go func() {
		<-ctx.Done()
		log.Printf("[hci_listener] shutting down HCI Listener")
		err := unix.Close(fd)
		if err != nil {
			return
		}
	}()

	//filter struct in kernel
	//struct hci_filter {
	//	uint32_t type_mask;
	//	uint32_t event_mask[2];
	//	uint16_t opcode;
	//};

	filter := []byte{
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff,
		0x00, 0x00,
		0x00, 0x00, // PADDING
	}

	err = unix.SetsockoptString(fd, unix.SOL_HCI, 2, string(filter))
	if err != nil {
		log.Println("[hci_listener] failed to set HCI filter", err)
		return "", err
	}

	//blueetooth chip
	addr := &unix.SockaddrHCI{
		Dev:     0, //hci0
		Channel: unix.HCI_CHANNEL_RAW,
	}
	if err := unix.Bind(fd, addr); err != nil {
		log.Println("[hci_listener] failed to bind HCI channel:", err)
		return "", err
	}

	log.Println("[hci_listener] blueetooth chip listening at the hardware level ")

	buf := make([]byte, 1024)

	for {
		n, err := unix.Read(fd, buf)
		if err != nil {
			log.Println("[hci_listener] failed to read HCI channel:", err)
			continue
		}

		packet := buf[:n]
		//0x04 packege type (event)
		if packet[0] == 0x04 {
			if packet[1] == 0x03 && len(packet) >= 12 {
				mac := fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X",
					packet[11], packet[10], packet[9], packet[8], packet[7], packet[6])
				log.Println("[hci_listener] target mac address", mac)

				err := guard.HandleMacEvent(mac)
				if err != nil {
					log.Println("[hci_listener] failed to handle mac:", err)
				}
			}
		}
	}
}

package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
)

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

func ResolveConfigPath() string {
	hasJson := FileExists("../config.json")
	hasConf := FileExists("../guard.conf")

	switch {
	case hasConf && hasJson:
		log.Println("WARNING: Both guard.conf and config.json found in the directory! Defaulting to guard.conf to prevent conflicts.")
		return "guard.conf"
	case hasConf:
		return "guard.conf"
	case hasJson:
		return "config.json"
	default:
		log.Println("WARNING: No configuration file found. Initializing guards with empty rules.")
		return ""
	}
}

func AskPermission(title, message string) bool {
	username := os.Getenv("SUDO_USER")

	if username == "" {
		username = os.Getenv("USER")
	}

	if username == "" {
		log.Println("user cannot found, cannot send notification")
		return false
	}

	userUid, err := user.Lookup(username)
	if err != nil {
		log.Println("Lookup error")
		return false
	}

	// cmd := exec.Command("sudo", "-u", username,
	// 	"env",
	// 	fmt.Sprintf("DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/%s/bus", userUid.Uid),
	// 	fmt.Sprintf("XDG_RUNTIME_DIR=/run/user/%s", userUid.Uid),
	// 	"DISPLAY=:0",
	// 	"notify-send", title, message, "--urgency=critical",
	// 	"--wait",
	// 	"--action=allow=Allow for this device",
	// )

	cmd := exec.Command("sudo", "-u", username,
		"env",
		fmt.Sprintf("DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/%s/bus", userUid.Uid),
		fmt.Sprintf("XDG_RUNTIME_DIR=/run/user/%s", userUid.Uid),
		"DISPLAY=:0",
		"zenity", "--question",
		"--title", title,
		"--text", message,
		"--width=350",
		"--ok-label=✅ Allow",
		"--cancel-label=❌ Ignore",
	)

	err = cmd.Run()

	return err == nil
}

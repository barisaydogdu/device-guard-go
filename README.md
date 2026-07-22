# Device Guard Go

Device Guard Go is a zero-dependency device security tool that protects your Linux system against unauthorized hardware access. It runs silently in the background with minimal resource consumption, instantly bringing USB, Network, and Bluetooth interfaces under control.
## Features

- It directly listens to the Linux kernel's uevent mechanism so as not to miss the moment USB and network devices are connected to the system.
- It monitors BlueZ signals in real time to capture Bluetooth device movements.
- The protection mechanism activates the moment a new device is connected to the system. The device is blocked until it is added to the safe list.
- When a new device is detected, the user is asked for confirmation. If the user grants permission (allow), the device is added to the whitelist; if they deny it, the device's connection to the system is immediately severed.

- It does not require any external third-party libraries or heavy dependencies. It uses only the Go standard library and native Linux APIs.


## Installation

Make sure Go is installed on your system.

```bash
 git clone https://github.com/barisaydogdu/device-guard-go.git
 cd device-guard-go
```
    
## Screenshots

<img width="932" height="790" alt="screenshot" src="https://github.com/user-attachments/assets/646b2a08-aa88-4166-a8f0-5c416e445119" />


## Tech Stack

* **Language:** Go (Golang) v1.26.3
* **Core APIs:** Linux Kernel Netlink (uevents), BlueZ D-Bus
## Contributing

Contributions are always welcome! 

Feel free to open an **Issue** to report bugs or suggest features, or submit a **Pull Request** directly with your improvements.

[![Go Reference](https://pkg.go.dev/badge/github.com/barisaydogdu/device-guard-go@v1.0.1#section-readme.svg)](https://pkg.go.dev/github.com/barisaydogdu/device-guard-go@v1.0.1#section-readme)


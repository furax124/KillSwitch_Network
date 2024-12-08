# KillSwitch_Network

Emergency KillSwitch is a program written in Go which turn off all network connections

## Features

- Toggle network connections using Alt+N hotkey (configurable)
- Disable all network adapters, WiFi, Ethernet, and network services
- Enable all network adapters, WiFi, Ethernet, and network services
- No pop up window using syscall package

## Requirements

- Go 1.23 or later
- Windows 11

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/KillSwitch_Network.git
   cd KillSwitch_Network
   go build -o KillSwitch_Network.exe -ldflags -H=windowsgui
   ```
   
## Utility

If a Stealer tries to steal your data, you can quickly disable all network connections using the Alt+N hotkey. This will prevent the Stealer from sending your data to their server.
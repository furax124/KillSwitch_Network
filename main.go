package main

import (
	"golang.design/x/hotkey"
	"golang.design/x/hotkey/mainthread"
	"log"
	"os/exec"
	"syscall"
	"time"
)

var networkEnabled = true

func main() {
	mainthread.Init(run)
}

func run() {
	// Register Alt+N hotkey
	hkAltN := hotkey.New([]hotkey.Modifier{hotkey.ModAlt}, hotkey.KeyN)

	if err := hkAltN.Register(); err != nil {
		log.Fatalf("[-] hotkey: failed to register hotkey Alt+N: %v", err)
		return
	}

	log.Printf("[+] Network Toggle Switch activated - Alt+N is registered\n")

	for {
		select {
		case <-hkAltN.Keydown():
			currentTime := time.Now().Format("2006-01-02 15:04:05")
			if networkEnabled {
				log.Printf("[*] %s - Kill switch activated, disabling all network connections...", currentTime)
				start := time.Now()
				disableAllNetworks()
				duration := time.Since(start)
				log.Printf("[+] Network disabled in %v\n", duration)
				networkEnabled = false
			} else {
				log.Printf("[*] %s - Enabling all network connections...", currentTime)
				start := time.Now()
				enableAllNetworks()
				duration := time.Since(start)
				log.Printf("[+] Network enabled in %v\n", duration)
				networkEnabled = true
			}
		}
	}
}

func disableAllNetworks() {
	commands := []struct {
		cmd  *exec.Cmd
		desc string
	}{
		{exec.Command("powershell", "-Command", `Get-NetAdapter | ForEach-Object { Disable-NetAdapter -Name $_.Name -Confirm:$false }`), "disable network adapters"},
		{exec.Command("netsh", "interface", "set", "interface", "Wi-Fi", "disable"), "disable WiFi"},
		{exec.Command("netsh", "interface", "set", "interface", "Ethernet", "disable"), "disable Ethernet"},
		{exec.Command("powershell", "-Command", `Get-WmiObject Win32_NetworkAdapter | Where-Object { $_.NetEnabled -eq $true } | ForEach-Object { $_.Disable() }`), "disable remaining interfaces"},
		{exec.Command("powershell", "-Command", `Stop-Service -Name "Wlan*" -Force; Stop-Service -Name "RmSvc" -Force; Stop-Service -Name "NetworkLocationWatcher" -Force; Stop-Service -Name "NlaSvc" -Force; Stop-Service -Name "netprofm" -Force`), "stop network services"},
	}

	for _, c := range commands {
		c.cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if err := c.cmd.Run(); err != nil {
			log.Printf("[-] Failed to %s: %v\n", c.desc, err)
		}
	}

	log.Println("[+] All network connections disabled successfully")
}

func enableAllNetworks() {
	commands := []struct {
		cmd  *exec.Cmd
		desc string
	}{
		{exec.Command("powershell", "-Command", `Get-NetAdapter | ForEach-Object { Enable-NetAdapter -Name $_.Name -Confirm:$false }`), "enable network adapters"},
		{exec.Command("netsh", "interface", "set", "interface", "Wi-Fi", "enable"), "enable WiFi"},
		{exec.Command("netsh", "interface", "set", "interface", "Ethernet", "enable"), "enable Ethernet"},
		{exec.Command("powershell", "-Command", `Get-WmiObject Win32_NetworkAdapter | ForEach-Object { $_.Enable() }`), "enable remaining interfaces"},
		{exec.Command("powershell", "-Command", `Start-Service -Name "Wlan*"; Start-Service -Name "RmSvc"; Start-Service -Name "NetworkLocationWatcher"; Start-Service -Name "NlaSvc"; Start-Service -Name "netprofm"`), "start network services"},
	}

	for _, c := range commands {
		c.cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if err := c.cmd.Run(); err != nil {
			log.Printf("Failed to %s: %v\n", c.desc, err)
		}
	}

	log.Println("[+] All network connections enabled successfully")
}

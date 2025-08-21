package main

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func main() {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Enum\USBSTOR`, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		fmt.Printf("Error opening registry key: %v\n", err)
		return
	}

	defer key.Close()

	// Read subkey names (USB storage devices)
	devices, err := key.ReadSubKeyNames(-1)
	if err != nil {
		fmt.Printf("Error reading subkey names: %v\n", err)
		return
	}

	// Iterate over device names
	for _, device := range devices {
		info, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Enum\USBSTOR\`+device, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
		if err != nil {
			fmt.Println("Erreur : ", err)
		}
		dateInfo, err := info.Stat()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("\nPERIPH : \n %s : %s \n\n", device, dateInfo.ModTime())
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"golang.org/x/sys/windows/registry"
)

func banner() {
	banner, _ := pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("USB", pterm.FgGray.ToStyle()),
		putils.LettersFromStringWithStyle("CHCK3R", pterm.FgGreen.ToStyle()),
	).Srender()
	fmt.Printf("\n")
	fmt.Println(banner)
}

var deviceVIDFound = false
var devicePID string
var deviceVID string

func fetchDevices(deviceVID, devicePID string) {
	file, err := os.OpenFile("usbdevices.txt", os.O_RDONLY, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, strings.ToLower(deviceVID)) {
			fmt.Println(line)
			deviceVIDFound = true
		}
		if deviceVIDFound && strings.HasPrefix(line, "\t") && strings.Contains(line, strings.ToLower(devicePID)) {
			fmt.Println(line)
		}
	}

}

func main() {
	banner()
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Enum\USB`, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
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
		_, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Enum\USB\`+device, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
		if err != nil {
			fmt.Println("Erreur : ", err)
		}
		// dateInfo, err := info.Stat()
		// if err != nil {
		// 	fmt.Println(err)
		// }

		deviceVIDtmp, after, _ := strings.Cut(device, "&")

		if strings.Contains(after, "&") {
			devicePIDtmp, _, _ := strings.Cut(after, "&")
			devicePID = devicePIDtmp
		} else {
			devicePID = after
		}

		deviceVID = strings.Trim(deviceVIDtmp, "VID_")
		devicePID = strings.Trim(devicePID, "PID_")
		// deviceMI := strings.Trim(deviceMItmp, "MI_")

		fetchDevices(deviceVID, devicePID)

		// fmt.Printf("%s : %s \n\n", device, dateInfo.ModTime().String())

	}
}

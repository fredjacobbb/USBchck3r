package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"golang.org/x/sys/windows/registry"
)

type USBDeviceInfo struct {
	Vendor []Vendor `json:"vendors"`
}

type Vendor struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Device    []Device `json:"devices"`
	InfoModif string
}

type Device struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

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
var deviceINFO string

func fetchDevices(deviceVID, devicePID string, deviceINFO string) {
	url := fmt.Sprintf("http://apps.sebastianlang.net/usb-ids?vid=%s&pid=%s", deviceVID, devicePID)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Request Error : ", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var usbDeviceInfo USBDeviceInfo

	err = json.Unmarshal(body, &usbDeviceInfo)
	if err != nil {
		return
	}

	fmt.Printf("----------\nUSBPERIPHERIK : %s\nDATEMODIF : %s\n----------", usbDeviceInfo.Vendor, deviceINFO)

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
		info, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Enum\USB\`+device, registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
		if err != nil {
			fmt.Println("Erreur : ", err)
		}

		defer info.Close()

		deviceINFO, err := info.Stat()
		if err != nil {
			return
		}

		deviceVIDtmp, after, _ := strings.Cut(device, "&")

		if strings.Contains(after, "&") {
			devicePIDtmp, _, _ := strings.Cut(after, "&")
			devicePID = devicePIDtmp
		} else {
			devicePID = after
		}

		deviceVID = strings.Trim(deviceVIDtmp, "VID_")
		devicePID = strings.Trim(devicePID, "PID_")
		fetchDevices(deviceVID, devicePID, deviceINFO.ModTime().Format("Monday, January 2, 2006 15:04:05 MST"))
		time.Sleep(time.Second * 2)

	}
}

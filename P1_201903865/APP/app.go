package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"os/exec"

	"github.com/pbnjay/memory"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) Cpu() float64 {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0
	} else {
		return percent[0]
	}
}

func (a *App) Disk() [2]uint64 {

	UsageStats, err := disk.Usage("/")
	var values [2]uint64

	if err != nil {
		values[0] = 0
		values[1] = 0
		return values
	}
	values[0] = UsageStats.Free
	values[1] = UsageStats.Used

	return values
}

func (a *App) Ram() [2]int {

	var values [2]int

	values[0] = int(memory.FreeMemory())
	values[1] = int(memory.TotalMemory() - memory.FreeMemory())

	return values
}

type USBport struct {
	VendorId  string `json:"vid"`
	ProductId string `json:"pid"`
	Name      string `json:"name"`
	Device    string `json:"device"`
	Status    bool   `json:"status"`
	Port_name string `json:"port_name"`
}

// Devolver nombres
func GetDeviceName(vendorId, productId string) (string, error) {
	devices, err := ioutil.ReadDir("/sys/bus/usb/devices")
	if err != nil {
		return "", err
	}

	for _, device := range devices {
		devicePath := filepath.Join("/sys/bus/usb/devices", device.Name())
		idVendor, err := ioutil.ReadFile(filepath.Join(devicePath, "idVendor"))
		if err != nil {
			continue
		}
		idProduct, err := ioutil.ReadFile(filepath.Join(devicePath, "idProduct"))
		if err != nil {
			continue
		}

		if strings.TrimSpace(string(idVendor)) == vendorId && strings.TrimSpace(string(idProduct)) == productId {
			return device.Name(), nil
		}
	}

	return "", errors.New("Dispositivo no encontrado")
}

// Listar Todas las USB
func (a *App) ListUSB() []USBport {

	cmd := exec.Command("lsusb")
	output, err := cmd.Output()

	if err != nil {
		fmt.Println("Error al ejecutar el comando: ", err)
		return nil
	}

	outputStr := string(output)
	outputLines := strings.Split(outputStr, "\n")

	var ports []USBport
	for _, line := range outputLines {
		if len(line) > 0 {
			fields := strings.Split(line, ":")

			port := USBport{
				VendorId:  strings.Fields(fields[1])[1],
				ProductId: strings.Fields(fields[2])[0],
				Device:    strings.Fields(fields[0])[3],
				Name:      strings.TrimLeft(fields[2], strings.Fields(fields[2])[0]),
			}

			devices, err := ioutil.ReadDir("/sys/bus/usb/devices")
			if err != nil {
				return nil
			}
			name, err := GetDeviceName(port.VendorId, port.ProductId)
			if err != nil {
				return nil
			}
			for _, device := range devices {
				if strings.Contains(device.Name(), name) {
					devicePath := filepath.Join("/sys/bus/usb/devices", device.Name())
					authorizedPath := filepath.Join(devicePath, "authorized")
					content, err := ioutil.ReadFile(authorizedPath)
					if err != nil {
						return nil
					}
					port.Port_name = device.Name()
					port.Status = strings.TrimSpace(string(content)) != "0"
				}
			}
			ports = append(ports, port)
		}
	}
	return ports
}

// Activar USB
// Desactivar USB
// Ver Logs

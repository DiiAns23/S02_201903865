package sys

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"bytes"
	"os/exec"
	"strings"

	"io/ioutil"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/pbnjay/memory"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/wailsapp/wails"
)

// Stats .
type Stats struct {
	log *wails.CustomLogger
}

// CPUUsage .
type CPUUsage struct {
	Average int `json:"avg"`
}

// DISKUsage
type DISKUsage struct {
	Used int `json:"used"`
	Free int `json:"free"`
}

// RAMUsage
type RAMUsage struct {
	Used int `json:"used"`
	Free int `json:"free"`
}

// USB
// USB Port
type USBport struct {
	VendorId  string `json:"vid"`
	ProductId string `json:"pid"`
	Name      string `json:"name"`
	Device    string `json:"device"`
	Status    bool   `json:"status"`
	Port_name string `json:"port_name"`
}

// USB Port Free
type USBport_free struct {
	Status    bool   `json:"status"`
	Port_name string `json:"port_name"`
}

// GetCPUUsage .
func (s *Stats) GetCPUUsage() *CPUUsage {
	percent, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		s.log.Errorf("No se puede obtener info del cpu: %s", err.Error())
		return nil
	}

	return &CPUUsage{
		Average: int(math.Round(percent[0])),
	}
}

// GetDISKUsage
func (s *Stats) GetDISKUsage() *DISKUsage {
	UsageStats, err := disk.Usage("/")
	if err != nil {
		s.log.Errorf("No se puede obtener info del disko: %s", err.Error())
		return nil
	}

	return &DISKUsage{
		Used: int(UsageStats.Used),
		Free: int(UsageStats.Free),
	}
}

func (s *Stats) GetRAMUsage() *RAMUsage {
	return &RAMUsage{
		Used: int(memory.TotalMemory() - memory.FreeMemory()),
		Free: int(memory.FreeMemory()),
	}
}

func (s *Stats) GetUSBPorts() []USBport {
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

func (s *Stats) DeshabilitarUSB(vendorId, productId, password string) (string, error) {
	devices, err := ioutil.ReadDir("/sys/bus/usb/devices")
	if err != nil {
		return "", err
	}

	name := vendorId
	if productId != "" {
		name2, err := GetDeviceName(vendorId, productId)
		if err != nil {
			return "", err
		}
		name = name2
	}

	for _, device := range devices {
		if strings.Contains(device.Name(), name) {
			devicePath := filepath.Join("/sys/bus/usb/devices", device.Name())
			authorizedPath := filepath.Join(devicePath, "authorized")
			cmd := exec.Command("sudo", "-S", "sh", "-c", fmt.Sprintf("echo 0 > %s", authorizedPath))
			cmd.Stdin = strings.NewReader(password + "\n")
			var out, stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil {
				return "", err
			}
			return "Dispositivo USB deshabilitado " + name, nil
		}
	}

	return "Error, dispositivo USB no existe", nil
}

func (s *Stats) HabilitarUSB(vendorId, productId, password string) (string, error) {
	devices, err := ioutil.ReadDir("/sys/bus/usb/devices")
	if err != nil {
		return "", err
	}

	name := vendorId
	if productId != "" {
		name2, err := GetDeviceName(vendorId, productId)
		if err != nil {
			return "", err
		}
		name = name2
	}

	for _, device := range devices {
		if strings.Contains(device.Name(), name) {
			devicePath := filepath.Join("/sys/bus/usb/devices", device.Name())
			authorizedPath := filepath.Join(devicePath, "authorized")
			cmd := exec.Command("sudo", "sh", "-c", fmt.Sprintf("echo 1 > %s", authorizedPath))
			cmd.Stdin = strings.NewReader(password + "\n")
			var out, stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil {
				return "", err
			}
			return "Dispositivo USB habilitado " + name, nil
		}
	}

	return "Error, dispositivo USB no existe", nil
}

func getMountedUSBPaths() ([]string, error) {
	var directorios []string
	directorio := "/media"
	users, _ := ioutil.ReadDir(directorio)
	for _, user := range users {
		user_path := directorio + "/" + user.Name()

		// Leer el contenido del directorio
		files, err := ioutil.ReadDir(user_path)
		if err != nil {
			return nil, err
		}

		// Filtrar solo los directorios
		for _, file := range files {
			if file.IsDir() {
				directorios = append(directorios, user_path+"/"+file.Name())
			}
		}
	}

	return directorios, nil
}

func (s *Stats) UsbLogs(password string) {
	// Crear un archivo .txt
	txtFile, err := os.OpenFile("bitacora.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer txtFile.Close()

	// Crear un watcher de inotify
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Obtener la lista de rutas de las unidades USB montadas
	usbPaths, err := getMountedUSBPaths()
	if err != nil {
		log.Fatal(err)
	}

	// Agregar todas las rutas de las unidades USB montadas al watcher
	for _, usbPath := range usbPaths {
		println("Path: " + usbPath)
		err = watcher.Add(usbPath)
		if err != nil {
			fmt.Printf("Error en watcher.Add: %v\n", err)
		}
	}

	fmt.Println("Esperando eventos de cambio en las unidades USB...")

	for {
		select {
		case event := <-watcher.Events:
			fmt.Println("Evento:", event)
			// Solo interesan los eventos de creación o eliminación de archivos
			if event.Op == fsnotify.Create || event.Op == fsnotify.Remove || event.Op == fsnotify.Rename || event.Op == fsnotify.Write {
				// Obtener la ruta de la unidad USB modificada
				usbPath := filepath.Dir(event.Name)

				// Leer los nombres de los archivos en la unidad USB modificada
				fileName := filepath.Base(event.Name)

				// Obtener la fecha y hora actual
				now := time.Now()

				// Escribir los nombres de los archivos modificados con la fecha y hora en una nueva línea en el archivo .txt
				_, err := fmt.Fprintf(txtFile, "[%s] - %s - %s\n", now.Format("2006-01-02 15:04:05"), event.Op, filepath.Join(usbPath, fileName))
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println("Se ha actualizado el archivo 'archivos_usb.txt' con los nombres de los archivos modificados en la unidad USB.")
			}
		case err := <-watcher.Errors:
			log.Println("Error:", err)
		}
	}
}

// Abrir el archivo .txt en una terminal con el comando nano
func (s *Stats) AbrirArchivo() {
	archivo := "bitacora.txt"

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fullPath := filepath.Join(workingDir, archivo)
	command := "nano " + fullPath

	cmd := exec.Command("gnome-terminal", "--", "bash", "-c", command)

	err2 := cmd.Run()
	if err2 != nil {
		fmt.Println("Error al abrir el archivo en la terminal:", err2)
		return
	}
}

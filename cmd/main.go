package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/DarkSoul94/deployer/menu"
	"github.com/DarkSoul94/deployer/template"
)

func main() {
	serviceData := menu.RunMenu()

	fileContent := template.CreateTemplate(serviceData)
	CreateFile(serviceData.Name, fileContent)

	OsExec(serviceData.Name)
}

func CreateFile(serviceName, serviceData string) {
	service_file, err := os.Create(fmt.Sprintf("/etc/systemd/system/%s.service", serviceName))
	if err != nil {
		panic(err)
	}
	defer service_file.Close()

	service_file.WriteString(serviceData)
}

func OsExec(serviceName string) {
	var err error

	reload := exec.Command("systemctl", "daemon-reload")
	err = reload.Run()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("- systemctl reload")
	}

	enable := exec.Command("systemctl", "enable", serviceName)
	err = enable.Run()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("- service enable")
	}

	start := exec.Command("systemctl", "start", serviceName)
	err = start.Run()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("- service start")
	}

	status := exec.Command("systemctl", "status", serviceName)
	stdout, err := status.Output()

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Print(string(stdout))
}

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/DarkSoul94/deployer/template"
)

func main() {
	var (
		service_name, filePath, wdPath string
		err                            error
	)

	service_name = ScanServiceName()
	filePath, wdPath = ScanPath()

	service_file, err := os.Create(fmt.Sprintf("/etc/systemd/system/%s.service", service_name))
	if err != nil {
		panic(err)
	}
	defer service_file.Close()

	service_file.WriteString(template.CreateTemplate(service_name, filePath, wdPath))

	reload := exec.Command("systemctl", "daemon-reload")
	err = reload.Run()
	if err != nil {
		fmt.Println(err)
	}

	start := exec.Command("systemctl", "start", service_name)
	err = start.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func ScanServiceName() string {
	var service_name string

name:
	fmt.Print("Введите имя сервиса: ")
	fmt.Fscan(os.Stdin, &service_name)

	if len(service_name) == 0 {
		fmt.Println("Имя сервиса не должно быть пустым")
		goto name
	}

	if _, err := os.Stat(fmt.Sprintf("/etc/systemd/system/%s.service", service_name)); err == nil {
		fmt.Println("Сервис с таким именем уже существует")
		goto name
	}

	return service_name
}

func ScanPath() (string, string) {
	var (
		filePath, wdPath string
	)

path:
	fmt.Print("Введите путь к исполняемому файлу: ")
	fmt.Fscan(os.Stdin, &filePath)
	if len(filePath) == 0 {
		fmt.Println("Путь не может быть пустым")
		goto path
	}

	if _, err := os.Stat(fmt.Sprintf("%s", filePath)); os.IsNotExist(err) {
		fmt.Println("Указаного файла не существует")
		goto path
	}

	s := strings.Split(filePath, "/")
	for i := 0; i < len(s)-1; i++ {
		wdPath += s[i] + "/"
	}

	return filePath, wdPath
}

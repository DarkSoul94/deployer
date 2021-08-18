package menu

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/DarkSoul94/deployer/models"
	"github.com/manifoldco/promptui"
)

func RunMenu() models.ServiceData {
	var data models.ServiceData
	data.Name = scanServiceName()
	data.PathToFile, data.PathToWorkingDirectory = scanPath()

	users := getUserList()
	data.User = selector(users, "Users")

	groups := getUserGroup(data.User)
	data.UserGroup = selector(groups, "Groups")

	return data
}

func scanServiceName() string {
	var serviceName string

name:
	fmt.Print("Введите имя сервиса: ")
	fmt.Fscan(os.Stdin, &serviceName)

	if len(serviceName) == 0 {
		fmt.Println("Имя сервиса не должно быть пустым")
		goto name
	}

	if _, err := os.Stat(fmt.Sprintf("/etc/systemd/system/%s.service", serviceName)); err == nil {
		fmt.Println("Сервис с таким именем уже существует")
		goto name
	}

	return serviceName
}

func scanPath() (string, string) {
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

func selector(slice []string, selectorName string) string {
	selector := promptui.Select{
		Label: selectorName,
		Items: slice,
	}

	_, result, err := selector.Run()
	if err != nil {
		fmt.Println(err)
	}

	return result
}

func getUserList() []string {
	getUsers := exec.Command("cut", "-d:", "-f1", "/etc/passwd")
	stdout, err := getUsers.Output()
	if err != nil {
		fmt.Println(err.Error())
	}

	users := strings.Split(string(stdout), "\n")

	return users[0 : len(users)-1]
}

func getUserGroup(user string) []string {
	getGroups := exec.Command("groups", user)
	stdout, err := getGroups.Output()
	if err != nil {
		fmt.Println(err.Error())
	}

	groups := strings.Split(
		strings.Split(string(stdout), " : ")[1],
		"\n",
	)

	return groups[0 : len(groups)-1]
}

package template

import (
	"strings"
	"text/template"

	"github.com/DarkSoul94/deployer/models"
)

type String string

func (s String) format(data map[string]interface{}) (out string, err error) {
	t := template.Must(template.New("").Parse(string(s)))
	builder := &strings.Builder{}
	if err = t.Execute(builder, data); err != nil {
		return
	}
	out = builder.String()
	return
}

func CreateTemplate(serviceData models.ServiceData) string {
	const template string = `
[Unit]
Description=a "{{.Name}}" service
After=syslog.target
After=network.target
	 
[Service]
Type=simple
Restart=on-failure
RestartSec=5
ExecStart={{.Path}}

User={{.User}}
Group={{.Group}}
		
WorkingDirectory={{.Directory}}
		
[Install]
WantedBy=multi-user.target`

	data := map[string]interface{}{
		"Name":      serviceData.Name,
		"Path":      serviceData.PathToFile,
		"Directory": serviceData.PathToWorkingDirectory,
		"User":      serviceData.User,
		"Group":     serviceData.UserGroup,
	}

	s, _ := String(template).format(data)

	return s
}

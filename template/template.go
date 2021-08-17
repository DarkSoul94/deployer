package template

import (
	"strings"
	"text/template"
)

type String string

func (s String) Format(data map[string]interface{}) (out string, err error) {
	t := template.Must(template.New("").Parse(string(s)))
	builder := &strings.Builder{}
	if err = t.Execute(builder, data); err != nil {
		return
	}
	out = builder.String()
	return
}

func CreateTemplate(name, pathToFile, workingDirectory string) string {
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

User=root
Group=root
		
WorkingDirectory={{.Directory}}
		
[Install]
WantedBy=multi-user.target`

	data := map[string]interface{}{
		"Name":      name,
		"Path":      pathToFile,
		"Directory": workingDirectory,
	}

	s, _ := String(template).Format(data)

	return s
}

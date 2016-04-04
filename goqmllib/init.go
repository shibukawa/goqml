package goqmllib

import (
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Init() {
	config, existing := LoadConfig()
	if !existing {
		config.Save()
	}

	color.Magenta("Creating main.go\n")
	ioutil.WriteFile("main.go", MustAsset("templates/main.go"), 0644)

	color.Magenta("Creating qtrc folder to store Qt resources\n")
	os.MkdirAll("qtrc", 0755)
	color.Magenta("Creating Resources folder to store binaries to bundle application\n")
	color.Magenta("Creating qtrc/main.qml\n")
	ioutil.WriteFile(filepath.Join("qtrc", "main.qml"), MustAsset("templates/main.qml"), 0644)

	os.MkdirAll("Resources", 0755)
	color.Magenta("Creating Requests/icon.png as a default icon\n")
	ioutil.WriteFile(filepath.Join("Resources", "icon.png"), MustAsset("resources/qt-logo.png"), 0644)
}

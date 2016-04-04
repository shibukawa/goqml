package goqmllib

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const packageFileName = "package.toml"

type PackageConfig struct {
	Name             string `toml:"name"`
	Description      string `toml:"description"`
	Author           string `toml:"author"`
	Organization     string `toml:"organization"`
	License          string `toml:"license"`
	Version          []int  `toml:"version"`
	ProjectStartYear int    `toml:"project_start_year"`
	AdditionalFrameworks []string `toml:"additional_frameworks"`
}

func (config *PackageConfig) Save() error {
	file, err := os.Create(packageFileName)
	if err != nil {
		return err
	}
	encoder := toml.NewEncoder(file)
	err = encoder.Encode(config)
	if err == nil {
		color.Magenta("Update: %s\n", packageFileName)
	} else {
		color.Red("Write file error: %s - %s\n", packageFileName, err.Error())
	}
	return err
}

func (config *PackageConfig) Copyright() string {
	year := time.Now().Year()
	startYear := config.ProjectStartYear
	if startYear == 0 {
		panic("start year is empty")
	}
	var names []string
	if config.Organization != "" {
		names = append(names, config.Organization)
	}
	if config.Author != "" {
		names = append(names, config.Author)
	}
	if year == startYear {
		return fmt.Sprintf("%d %s", year, strings.Join(names, " "))
	}
	return fmt.Sprintf("%d-%d %s", startYear, year, strings.Join(names, " "))
}

func UserName() string {
	names := []string{"LOGNAME", "USER", "USERNAME"}
	for _, name := range names {
		value := os.Getenv(name)
		if value != "" {
			return value
		}
	}
	return "(no name)"
}

func LoadConfig() (*PackageConfig, bool) {
	file, err := os.Open(packageFileName)
	if err == nil {
		config := &PackageConfig{}
		_, err := toml.DecodeReader(file, config)
		if err == nil {
			return config, true
		}
	}
	abs, _ := filepath.Abs(".")
	name := filepath.Base(abs)

	return &PackageConfig{
		Name:             strings.ToUpper(name[:1]) + name[1:],
		Description:      "write description here",
		Author:           UserName(),
		Organization:     "Your Compony or Organization",
		License:          "BSD",
		Version:          []int{1, 0, 0},
		ProjectStartYear: time.Now().Year(),
		AdditionalFrameworks: []string{"QtWebEngine"},
	}, false
}

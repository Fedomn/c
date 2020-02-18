package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

var configFile string

func init() {
	configFile = filepath.Dir(os.Args[0]) + "/.c.conf"
}

type Cmd struct {
	Cmd   string `yaml:"cmd"`
	Name  string `yaml:"name"`
	Alias string `yaml:"alias"`
}

func LoadCommands() []Cmd {
	var commands []Cmd
	var data []byte
	var err error
	if data, err = ioutil.ReadFile(configFile); err != nil {
		color.Green("Init bootstrap demo commands, please modify it: %s", configFile)
		return initBootstrapCommands()
	}
	if err = yaml.Unmarshal(data, &commands); err != nil {
		color.Red("Failed to parse %s", configFile)
		os.Exit(1)
	}
	return commands
}

func initBootstrapCommands() []Cmd {
	data := []byte(`-
 name: show ip
 cmd: curl https://ifconfig.co/json`)
	var commands []Cmd
	if err := yaml.Unmarshal(data, &commands); err != nil {
		color.Red("Init bootstrap commands failed")
		os.Exit(1)
	}
	if err := ioutil.WriteFile(configFile, data, 0644); err != nil {
		color.Red("Init bootstrap commands failed")
		os.Exit(1)
	}
	return commands
}

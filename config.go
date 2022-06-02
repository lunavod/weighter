package main

import (
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

type Scales struct {
	IP             string
	Port           int
	COMPort        string
	ConnectionType string
}

type Redis struct {
	Addr     string
	Password string
	Db       int
}

type Config struct {
	Scales Scales
	Redis  Redis
}

func GetConfig() Config {
	dat, err := ioutil.ReadFile("config.toml")
	if err != nil {
		panic("Config file not found")
	}

	config := Config{}
	toml.Unmarshal(dat, &config)

	if config.Scales.IP != "" && config.Scales.COMPort != "" && config.Scales.ConnectionType == "" {
		panic("Connection type not specified")
	}

	if config.Scales.IP != "" && config.Scales.ConnectionType == "" {
		config.Scales.ConnectionType = "IP"
	}

	if config.Scales.COMPort != "" && config.Scales.ConnectionType == "" {
		config.Scales.ConnectionType = "COM"
	}

	if config.Scales.COMPort == "" && config.Scales.IP == "" {
		panic("Connection params not specified")
	}

	return config
}

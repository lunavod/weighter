package main

import (
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

type Scales struct {
	IP string
	Port int
	COMPort string
	ConnectionType string
}

type Server struct {
	IP string
	Port int
}

type Config struct {
	Scales Scales
	Server Server
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

	if config.Server.IP == "" || config.Server.Port == 0 {
		panic("Server ip or port not specified")
	}

	return config
}

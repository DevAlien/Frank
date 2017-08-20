package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"frank/src/go/models"
)

type Config struct {
	Name         string                 `json:"name"`
	Deactivation []string               `json:"deactivation"`
	Configs      map[string]interface{} `json:"configs"`
	Devices      []Device               `json:"devices"`
	NamedDevices map[string]Device
	Commands     []Command `json:"commands"`
}

type Device struct {
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	Connection DeviceConnection  `json:"connection"`
	Interfaces []DeviceInterface `json:"interfaces"`
}

type DeviceConnection struct {
	Type    string `json:"type"`
	Address string `json:"address"`
}

type DeviceInterface struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
	Pin    string `json:"pin"`
}

type Command struct {
	Name     string                 `json:"name"`
	Commands []string               `json:"commands"`
	Actions  []models.CommandAction `json:"actions"`
}

type MatchingInterface struct {
}

var ParsedConfig Config

func (c *Config) Get(key string) string {
	if val, ok := c.Configs[key]; ok {
		return val.(string)
	}

	return ""
}

func GetConfig(fileName string) (Config, error) {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(fileName)
		log.Fatal(err)
	}

	return ParseConfig(content)
}

func ParseConfig(input []byte) (Config, error) {

	if err := json.Unmarshal(input, &ParsedConfig); err != nil {
		return ParsedConfig, err
	}

	ParsedConfig.NamedDevices = map[string]Device{}

	if len(ParsedConfig.Devices) > 0 {
		for _, device := range ParsedConfig.Devices {
			ParsedConfig.NamedDevices[device.Name] = device
		}
	}

	return ParsedConfig, nil
}

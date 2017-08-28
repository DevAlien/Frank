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
	NamedDevices map[string]Device      `json:"-"`
	Commands     []Command              `json:"commands"`
}

type Device struct {
	Name       string            `json:"name" binding:"required"`
	Type       string            `json:"type" binding:"required"`
	Connection DeviceConnection  `json:"connection"`
	Interfaces []DeviceInterface `json:"interfaces"`
	Commands   []Command         `json:"commands,omitempty"`
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

var FileName string

func (c *Config) Get(key string) string {
	if val, ok := c.Configs[key]; ok {
		return val.(string)
	}

	return ""
}

func InitConfig(fileName string) error {
	FileName = fileName
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(fileName)
		log.Fatal(err)
	}

	return parseConfig(content)
}

func saveConfig() error {
	err := ioutil.WriteFile(FileName, toJSON(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func toJSON() []byte {
	b, _ := json.MarshalIndent(ParsedConfig, "", "  ")

	return b
}

func AddCommand(command Command) error {
	ParsedConfig.Commands = append(ParsedConfig.Commands, command)
	return saveConfig()
}

func AddDevice(device Device) error {
	ParsedConfig.Devices = append(ParsedConfig.Devices, device)
	generateParsedDevices()

	return saveConfig()
}

func RemoveDevice(deviceName string) error {
	found := false
	devices := []Device{}
	for _, d := range ParsedConfig.Devices {
		if d.Name != deviceName {
			devices = append(devices, d)
		} else {
			found = true
		}
	}

	if found == false {
		return fmt.Errorf("The device %s was not found, therefore was not delete", deviceName)
	}

	ParsedConfig.Devices = devices
	return nil
}

func GetCommandsByDeviceName(deviceName string) []Command {
	commands := []Command{}

	for _, c := range ParsedConfig.Commands {
		for _, cA := range c.Actions {
			if cA.DeviceName == deviceName {
				commands = append(commands, c)
				break
			}
		}
	}

	return commands
}

func GetDevice(deviceName string) (Device, error) {
	device := Device{}
	if device, ok := ParsedConfig.NamedDevices[deviceName]; ok {
		return device, nil
	}

	return device, fmt.Errorf("Device \"%s\" not found", deviceName)
}

func GetDeviceInterface(deviceName string, interfaceName string) (DeviceInterface, error) {
	deviceInterface := DeviceInterface{}
	device, err := GetDevice(deviceName)
	if err != nil {
		return deviceInterface, err
	}

	return getDeviceInterface(device, interfaceName)
}

func getDeviceInterface(device Device, interfaceName string) (DeviceInterface, error) {
	for _, in := range device.Interfaces {
		if in.Name == interfaceName {
			return in, nil
		}
	}

	return DeviceInterface{}, fmt.Errorf("Interface \"%s\" not found in Device \"%s\"", interfaceName, device.Name)
}

func generateParsedDevices() {
	ParsedConfig.NamedDevices = map[string]Device{}

	if len(ParsedConfig.Devices) > 0 {
		for _, device := range ParsedConfig.Devices {
			ParsedConfig.NamedDevices[device.Name] = device
		}
	}
}

func parseConfig(input []byte) error {

	if err := json.Unmarshal(input, &ParsedConfig); err != nil {
		return err
	}

	generateParsedDevices()

	return nil
}

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"frank/src/go/models"
)

type Config struct{}

var ParsedConfig models.Config

var FileName string

func Get(key string) string {
	if val, ok := ParsedConfig.Configs[key]; ok {
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

func AddCommand(command models.Command) error {
	ParsedConfig.Commands = append(ParsedConfig.Commands, command)
	return saveConfig()
}

func AddDevice(device models.Device) error {
	ParsedConfig.Devices = append(ParsedConfig.Devices, device)
	generateParsedDevices()

	return saveConfig()
}

func AddAction(action models.Action) error {
	ParsedConfig.Actions = append(ParsedConfig.Actions, action)
	generateParsedActions()

	return saveConfig()
}

func RemoveDevice(deviceName string) error {
	found := false
	devices := []models.Device{}
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

func GetActionsByDeviceName(deviceName string) []models.Action {
	actions := []models.Action{}

	for _, a := range ParsedConfig.Actions {
		if a.DeviceName == deviceName {
			actions = append(actions, a)
		}
	}
	fmt.Println("actions", actions)
	return actions
}

func GetCommandsByDeviceName(deviceName string) []models.Command {
	commands := []models.Command{}

	for _, c := range ParsedConfig.Commands {
		for _, cA := range c.Actions {
			fmt.Println(cA)
			// if cA.DeviceName == deviceName {
			// 	commands = append(commands, c)
			// 	break
			// }
		}
	}

	return commands
}

func GetAvailablePlugins() []models.Plugin {
	plugins := []models.Plugin{}
	mappedPlugins := map[string]*models.Plugin{}

	for _, a := range ParsedConfig.Actions {
		if a.Plugin != "" {
			if _, ok := mappedPlugins[a.Plugin]; !ok {
				mappedPlugins[a.Plugin] = &models.Plugin{
					Name:    a.Plugin,
					Actions: []models.Action{},
				}
			}
			mappedPlugins[a.Plugin].Actions = append(mappedPlugins[a.Plugin].Actions, a)
		}
	}

	for _, v := range mappedPlugins {
		plugins = append(plugins, *v)
	}

	return plugins
}
func GetDevice(deviceName string) (models.Device, error) {
	device := models.Device{}
	if device, ok := ParsedConfig.NamedDevices[deviceName]; ok {
		return device, nil
	}

	return device, fmt.Errorf("Device \"%s\" not found", deviceName)
}

func GetAction(actionName string) (models.Action, error) {
	action := models.Action{}
	if action, ok := ParsedConfig.NamedActions[actionName]; ok {
		return action, nil
	}

	return action, fmt.Errorf("Action \"%s\" not found", actionName)
}

func GetDeviceInterface(deviceName string, interfaceName string) (models.DeviceInterface, error) {
	deviceInterface := models.DeviceInterface{}
	device, err := GetDevice(deviceName)
	if err != nil {
		return deviceInterface, err
	}

	return getDeviceInterface(device, interfaceName)
}

func getDeviceInterface(device models.Device, interfaceName string) (models.DeviceInterface, error) {
	for _, in := range device.Interfaces {
		if in.Name == interfaceName {
			return in, nil
		}
	}

	return models.DeviceInterface{}, fmt.Errorf("Interface \"%s\" not found in Device \"%s\"", interfaceName, device.Name)
}

func generateParsedDevices() {
	ParsedConfig.NamedDevices = map[string]models.Device{}

	if len(ParsedConfig.Devices) > 0 {
		for _, device := range ParsedConfig.Devices {
			ParsedConfig.NamedDevices[device.Name] = device
		}
	}
}

func generateParsedActions() {
	ParsedConfig.NamedActions = map[string]models.Action{}

	if len(ParsedConfig.Actions) > 0 {
		for _, action := range ParsedConfig.Actions {
			ParsedConfig.NamedActions[action.Name] = action
		}
	}
}

func parseConfig(input []byte) error {

	if err := json.Unmarshal(input, &ParsedConfig); err != nil {
		return err
	}

	generateParsedDevices()
	generateParsedActions()

	return nil
}

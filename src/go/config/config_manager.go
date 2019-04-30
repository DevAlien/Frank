package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	l "frank/src/go/helpers/log"
	"frank/src/go/models"

	"github.com/creasty/defaults"
	"github.com/radovskyb/watcher"
	"time"
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

func Is(key string) bool {
	if val, ok := ParsedConfig.Configs[key]; ok {
		return val.(bool)
	}

	return false
}

func readFile(fileName string) []byte {
	FileName = fileName
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(fileName)
		log.Fatal(err)
	}

	return content
}

func startWatcher(fileName string) {
	w := watcher.New()
	w.SetMaxEvents(1)

	w.FilterOps(watcher.Write)

	go func() {
		for {
			select {
			case <-w.Event:
				l.Log.Info("Config File Changed")
				parseConfig(readFile(fileName))
			case err := <-w.Error:
				l.Log.Error(err.Error())
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	if err := w.Add(fileName); err != nil {
		log.Fatalln(err)
	}

	// Start the watching process - it'll check for changes every 100ms.
	go func() {
		if err := w.Start(time.Millisecond * 100); err != nil {
			l.Log.Error(err.Error())
		}
	}()
}

func InitConfig(fileName string) error {
	content := readFile(fileName)
	startWatcher(fileName)

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

func SetDdns(ddns models.Ddns) error {
	ParsedConfig.Ddns = ddns
	return saveConfig()
}

func AddCommand(command models.Command) error {
	ParsedConfig.Commands = append(ParsedConfig.Commands, command)
	return saveConfig()
}

func AddDevice(device models.Device) error {
	ParsedConfig.Devices = append(ParsedConfig.Devices, device)

	return saveConfig()
}

func AddAction(action models.Action) error {
	ParsedConfig.Actions = append(ParsedConfig.Actions, action)

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

func GetDdns() models.Ddns {
	return ParsedConfig.Ddns
}

func GetVoice() *models.Voice {
	return &ParsedConfig.Voice
}

func GetHTTP() *models.HTTP {
	return &ParsedConfig.HTTP
}

func GetWebSocket() *models.WebSocket {
	return &ParsedConfig.WebSocket
}

func GetTelegram() *models.Telegram {
	return &ParsedConfig.Telegram
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

func GetReading(readingName string) (models.Reading, error) {
	reading := models.Reading{}

	if reading, ok := ParsedConfig.NamedReadings[readingName]; ok {
		return reading, nil
	}

	return reading, fmt.Errorf("Reading \"%s\" not found", readingName)
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

func generateParsedReadings() {
	ParsedConfig.NamedReadings = map[string]models.Reading{}

	if len(ParsedConfig.Readings) > 0 {
		for _, reading := range ParsedConfig.Readings {
			ParsedConfig.NamedReadings[reading.Name] = reading
		}
	}
}

func parseConfig(input []byte) error {
	l.Log.Debug("Parsing Config")
	ParsedConfig = models.Config{}
	if err := defaults.Set(&ParsedConfig); err != nil {
		return err
	}
	if err := json.Unmarshal(input, &ParsedConfig); err != nil {
		return err
	}
	fmt.Printf("%+v\n", ParsedConfig)
	generateParsedDevices()
	generateParsedActions()
	generateParsedReadings()

	return nil
}

package managers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	Name     string          `json:"name"`
	Commands []string        `json:"commands"`
	Actions  []CommandAction `json:"actions"`
}

type CommandAction struct {
	DeviceName        string                 `json:"device"`
	MatchingInterface map[string]interface{} `json:"matchingInterface"`
	InterfaceName     string                 `json:"interface"`
	Action            string                 `json:"action"`
}

type MatchingInterface struct {
}

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
	var config Config

	if err := json.Unmarshal(input, &config); err != nil {
		return config, err
	}

	config.NamedDevices = map[string]Device{}

	if len(config.Devices) > 0 {
		for _, device := range config.Devices {
			config.NamedDevices[device.Name] = device
		}
	}

	fmt.Printf("%+v\n", config.NamedDevices)
	return config, nil
}

// func main() {
// 	// our target will be of type map[string]interface{}, which is a pretty generic type
// 	// that will give us a hashtable whose keys are strings, and whose values are of
// 	// type interface{}
// 	var val map[string]interface{}
// 	var config ConfigManager
// 	cM := ConfigManager{}

// 	if err := json.Unmarshal([]byte(input), &config); err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("%+v\n", config)
// 	if err := json.Unmarshal([]byte(input), &val); err != nil {
// 		panic(err)
// 	}

// 	fmt.Println(val)
// 	cM.Configs = val["configs"].(map[string]interface{})
// 	fmt.Println(cM.Configs["google_api_key"].(string))
// 	for k, v := range val {

// 		fmt.Println(k, reflect.TypeOf(v))
// 	}
// }

package models

type Config struct {
	Name         string                 `json:"name"`
	Deactivation []string               `json:"deactivation"`
	Configs      map[string]interface{} `json:"configs"`
	Devices      []Device               `json:"devices"`
	NamedDevices map[string]Device      `json:"-"`
	NamedActions map[string]Action      `json:"-"`
	Actions      []Action               `json:"actions"`
	Commands     []Command              `json:"commands"`
}

type Device struct {
	Name       string            `json:"name" binding:"required"`
	Type       string            `json:"type" binding:"required"`
	Connection DeviceConnection  `json:"connection"`
	Interfaces []DeviceInterface `json:"interfaces"`
	Commands   []Command         `json:"commands,omitempty"`
	Actions    []Action          `json:"actions,omitempty"`
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
	Action         string                 `json:"action,omitempty"`
	MatchingAction map[string]interface{} `json:"matchingAction,omitempty"`
}

type Action struct {
	Name              string                 `json:"name"`
	DeviceName        string                 `json:"device,omitempty"`
	MatchingInterface map[string]interface{} `json:"matchingInterface,omitempty"`
	InterfaceName     string                 `json:"interface,omitempty"`
	Plugin            string                 `json:"plugin,omitempty"`
	Action            map[string]interface{} `json:"action"`
}

type Plugin struct {
	Name    string   `json:"name"`
	Actions []Action `json:"actions,omitempty"`
}

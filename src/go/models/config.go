package models

// Config is the struct for the config.json file
type Config struct {
	Name          string                 `json:"name"`
	Deactivation  []string               `json:"deactivation"`
	Configs       map[string]interface{} `json:"configs"`
	Devices       []Device               `json:"devices"`
	NamedDevices  map[string]Device      `json:"-"`
	NamedActions  map[string]Action      `json:"-"`
	NamedReadings map[string]Reading     `json:"-"`
	Readings      []Reading              `json:"readings"`
	Actions       []Action               `json:"actions"`
	Commands      []Command              `json:"commands"`
	Crons         []Cron                 `json:"crons"`
	Ddns          Ddns                   `json:"ddns,omitempty"`
	Voice         Voice                  `json:"voice,omitempty"`
	Telegram      Telegram               `json:"telegram,omitempty"`
	HTTP          HTTP                   `json:"http,omitempty"`
	WebSocket     WebSocket              `json:"web_socket,omitempty"`
}

// MqttLocal struct to define the structure of a local Mqtt server inside the Mqtt definition
type MqttLocal struct {
	Port     int  `json:"port"`
	Disabled bool `json:"disabled"`
}

// Mqtt struct to define the structure of the config json for Mqtt
type Mqtt struct {
	Local    MqttLocal `json:"local"`
	Address  string    `json:"address"`
	Port     int       `json:"port"`
	Username string    `json:"username"`
	Password string    `json:"Password"`
	Disabled bool      `json:"disabled"`
}

// WebSocket struct to define the structure of the websocket server
// You can define a Port or disable it
type WebSocket struct {
	Port     int  `json:"port"`
	Disabled bool `json:"disabled"`
}

// HTTP struct to define the structure of the HTTP server
// You can define a Port or disable it
type HTTP struct {
	Port     int  `json:"port"`
	Disabled bool `json:"disabled"`
}

// Voice struct to define the structure of the Voice service
// You can define an API key, a type or disable it
type Voice struct {
	APIKey   string `json:"api_key"`
	Type     string `json:"type"`
	Disabled bool   `json:"disabled"`
}

// Telegram struct to define the structure of the Telegram service
// You can define an API key or disable it
type Telegram struct {
	APIKey   string `json:"api_key"`
	Disabled bool   `json:"disabled"`
}

// Cron struct to define the structure of a Cron
type Cron struct {
	Description    string            `json:"description" binding:"required"`
	Action         string            `json:"action" binding:"required"`
	Every          int               `json:"every"`
	TimeType       string            `json:"time_type"`
	At             string            `json:"at,omitempty"`
	Extra          map[string]string `json:"extra,omitempty"`
	CronExpression string            `json:"cron_expression"`
}

// Device struct to define the structure of a Device
type Device struct {
	Name       string            `json:"name" binding:"required"`
	Type       string            `json:"type" binding:"required"`
	Connection DeviceConnection  `json:"connection"`
	Interfaces []DeviceInterface `json:"interfaces"`
	Commands   []Command         `json:"commands,omitempty"`
	Actions    []Action          `json:"actions,omitempty"`
}

// Ddns struct to define the structure of the Ddns service
type Ddns struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	Type           string `json:"type"`
	Hostname       string `json:"hostname"`
	CronExpression string `json:"cron_expression"`
}

// DeviceConnection struct to define the structure of the a Device Connection
type DeviceConnection struct {
	Type    string `json:"type"`
	Address string `json:"address"`
}

// DeviceInterface struct to define the structure of an interface for a device
// you define a driver, a pin and a name
type DeviceInterface struct {
	Name   string `json:"name"`
	Driver string `json:"driver"`
	Pin    string `json:"pin"`
}

// Command struct to define the structure of a Command
type Command struct {
	Name     string          `json:"name"`
	Commands []string        `json:"commands"`
	Actions  []CommandAction `json:"actions"`
}

// CommandAction struct to define the structure of a Command Action which
type CommandAction struct {
	Action         string                 `json:"action,omitempty"`
	MatchingAction map[string]interface{} `json:"matchingAction,omitempty"`
}

// Reading struct to define the structure of an Reading, which executes something on a device/plugin
type Reading struct {
	Name       string                 `json:"name"`
	DeviceName string                 `json:"device,omitempty"`
	Plugin     string                 `json:"plugin,omitempty"`
	Type       string                 `json:"type"`
	Params     map[string]interface{} `json:"params"`
}

// Action struct to define the structure of an Action, which executes something on a device/plugin
type Action struct {
	Name              string                 `json:"name"`
	DeviceName        string                 `json:"device,omitempty"`
	MatchingInterface map[string]interface{} `json:"matchingInterface,omitempty"`
	InterfaceName     string                 `json:"interface,omitempty"`
	Plugin            string                 `json:"plugin,omitempty"`
	Action            map[string]interface{} `json:"action"`
}

// Plugin struct to define the structure of a Plugin
type Plugin struct {
	Name    string   `json:"name"`
	Actions []Action `json:"actions,omitempty"`
}

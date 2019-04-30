package services

import (
	"frank/src/go/config"
)

//LoadServices loads the services needed to run the system based on the specified configuration.
func LoadServices() {
	if config.ParsedConfig.Ddns.Hostname != "" {
		LoadDdns(config.ParsedConfig.Ddns)
	}

	if config.GetVoice().APIKey != "" && config.GetVoice().Disabled != true {
		NewVoiceRecognition(config.GetVoice())
	}

	if config.GetWebSocket().Disabled != true {
		NewSocketIoServer(config.GetWebSocket())
	}

	if config.GetHTTP().Disabled != true {
		NewHTTPServer(config.GetHTTP())
	}

	if config.GetHTTP().Disabled != true {
		NewMqttServer(config.GetHTTP())
		NewMqttClient(config.GetHTTP())
	}

	if config.GetTelegram().Disabled != true {
		NewTelegramBotService(config.GetTelegram())
	}
}

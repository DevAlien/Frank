package devices

import (
	"frank/src/go/helpers/log"
	"frank/src/go/managers"
	"frank/src/go/services"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

type Devices struct {

}

func ManageCommands(commandsFound []services.CommandFound, config *managers.Config) {
	if len(commandsFound) < 1 {
		return
	}

	for _, co := range commandsFound {
		log.Log.Debug("Handling command", co.CommandName)
		for _, a := range co.Actions {
			HandleActions(a, co.ExtraText, config)
		}
	}
}

func HandleActions(action managers.CommandAction, extraText map[string]string, config *managers.Config) {
	device := config.NamedDevices[action.DeviceName]

	switch device.Type {
	case "firmata":
		FirmataHandler(action, device, extraText, config)
	}
}

func FirmataHandler(action managers.CommandAction, device managers.Device, extraText map[string]string, config *managers.Config) {
	var firmataA *firmata.Adaptor
	if device.Connection.Type == "tcp" {
		firmataAdaptor := firmata.NewTCPAdaptor(device.Connection.Address)
		firmataA = firmataAdaptor.Adaptor
	} else {
		firmataAdaptor := firmata.NewAdaptor(device.Connection.Address)
		firmataA = firmataAdaptor
	}

	err := firmataA.Connect()
	if err != nil {
		log.Log.Error("Could not connect to", device.Name, "at", device.Connection.Address)
		return
	}

	for _, di := range device.Interfaces {
		if len(action.MatchingInterface) > 0 {
			colore := action.MatchingInterface["colore"].( map[string]interface {})
			if val, ok := colore[extraText["colore"]].(string); ok {
				if val == di.Name {
					switch di.Driver {
						case "led":
							FirmataLedInterface(firmataA, action, di)
					}
				}
			}
			
		} else 
		if di.Name == action.InterfaceName {
			switch di.Driver {
				case "led":
					FirmataLedInterface(firmataA, action, di)
			}
		}
	}
	err = firmataA.Disconnect()
	if err != nil {
		log.Log.Notice("Could not disconnect from", device.Name)
		return
	}
}

func FirmataLedInterface(firmataA *firmata.Adaptor, action managers.CommandAction, deviceInterface managers.DeviceInterface) {
	log.Log.Debug("Interacting with led", deviceInterface.Pin)
	led := gpio.NewLedDriver(firmataA, deviceInterface.Pin)
	if action.Action == "on" {
		led.On()
	} else if action.Action == "off" {
		led.Off()
	} else {
		led.Toggle()
	}
}
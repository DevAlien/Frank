package plugins

import (
	"frank/src/go/config"
	"frank/src/go/helpers/log"
	"frank/src/go/models"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

type PluginFirmata struct {
}

func NewPluginFirmata() PluginFirmata {
	pf := PluginFirmata{}

	return pf
}

func (ctx *PluginFirmata) ExecAction(action models.Action, extraText map[string]string) {
	device, err := config.GetDevice(action.DeviceName)
	if err != nil {
		log.Log.Error(err.Error())
		return
	}
	log.Log.Debugf("Interacting with device %s", device.Name)
	go FirmataHandler(action, device, extraText)
}

func FirmataHandler(action models.Action, device models.Device, extraText map[string]string) {
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
			colore := action.MatchingInterface["colore"].(map[string]interface{})
			if val, ok := colore[extraText["colore"]].(string); ok {
				if val == di.Name {
					switch di.Driver {
					case "led":
						FirmataLedInterface(firmataA, action, di)
					}
				}
			}

		} else if di.Name == action.InterfaceName {
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

func FirmataLedInterface(firmataA *firmata.Adaptor, action models.Action, deviceInterface models.DeviceInterface) {
	log.Log.Debug("Interacting with led", deviceInterface.Pin)
	led := gpio.NewLedDriver(firmataA, deviceInterface.Pin)
	if action.Action["action"].(string) == "on" {
		led.On()
	} else if action.Action["action"].(string) == "off" {
		led.Off()
	} else {
		led.Toggle()
	}
}

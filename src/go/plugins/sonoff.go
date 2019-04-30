package plugins

import (
	"frank/src/go/config"
	"frank/src/go/helpers/log"
	"frank/src/go/models"
)

type PluginSonoff struct {
}

func NewPluginSonoff() PluginSonoff {
	ps := PluginSonoff{}

	return ps
}

func (ctx *PluginSonoff) ExecAction(action models.Action, extraText map[string]string) {
	device, err := config.GetDevice(action.DeviceName)
	if err != nil {
		log.Log.Error(err.Error())
		return
	}
	log.Log.Debugf("Interacting with device %s", device.Name)

	go SonoffHandler(action, device, extraText)
}

func SonoffHandler(action models.Action, device models.Device, extraText map[string]string) {

}

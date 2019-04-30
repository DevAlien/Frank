package plugins

import (
	"errors"
	"fmt"

	"frank/src/go/config"
	"frank/src/go/helpers/log"
	"frank/src/go/models"

	"github.com/imroc/req"
	"strings"
)

//PluginHTTP struct definition.
type PluginHTTP struct {
}

//NewPluginHTTP instantiate a new PluginHTTP.
func NewPluginHTTP() PluginHTTP {
	pf := PluginHTTP{}

	return pf
}

//ExecReading is from the Plugin interface and it is needed to execute a reading definied in the config.
//The PluginHTTP ExecReading gets the url from the params from the reading and does an HTTP request.
func (ctx *PluginHTTP) ExecReading(reading models.Reading) (models.ReadingResponse, error) {
	params := reading.Params
	response := models.ReadingResponse{
		Reading: reading,
	}

	if val, ok := params["url"]; ok {
		log.Log.Debugf("Calling the URL %s", val.(string))
		resp, err := req.Get(val.(string))
		if err != nil {
			log.Log.Error(err.Error())
			return response, err
		}
		r := resp.Response()
		log.Log.Info("Response, %s", string(resp.Bytes()))
		fmt.Printf("%v", strings.Contains(r.Header.Get("Content-Type"), "json"))
		response.Data = string(resp.Bytes())

		return response, nil
	}

	return response, errors.New("Could not find url parameter")
}

//ExecAction is from the Plugin interface and it is needed to execute an action definied in the config.
//The PluginHTTP ExecAction gets the URL from the action and executes the HTTP call.
func (ctx *PluginHTTP) ExecAction(action models.Action, extraText map[string]string) {
	if action.DeviceName == "" {
		log.Log.Debug("no device")
		go httpHandler(action, extraText)
		return
	}

	device, err := config.GetDevice(action.DeviceName)
	if err != nil {
		log.Log.Error(err.Error())
		return
	}
	log.Log.Debugf("Interacting with device '%s' doing '%s'", device.Name, action.Name)
	go httpHandlerDevice(action, device, extraText)
}

func httpHandler(action models.Action, extraText map[string]string) {
	log.Log.Debug(fmt.Sprintf("%s", action.Action["action"]))
	_, err := req.Get(fmt.Sprintf("%s", action.Action["action"]))
	if err != nil {
		log.Log.Error(err.Error())
	}
}

func httpHandlerDevice(action models.Action, device models.Device, extraText map[string]string) {
	url := strings.Replace(device.Connection.Address, "%", "%%", -1)
	log.Log.Debug(fmt.Sprintf("%+v", action.Action))
	log.Log.Debug(fmt.Sprintf("%+v", action.Action["action"]))
	log.Log.Debug(fmt.Sprintf("%s%s", url, action.Action["action"]))
	r, err := req.Get(fmt.Sprintf("%s%s", url, action.Action["action"]))
	if err != nil {
		log.Log.Error(err.Error())
	}

	log.Log.Debug(r.ToString())
}

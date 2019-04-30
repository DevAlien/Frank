package managers

import (
	"errors"

	"frank/src/go/config"
	"frank/src/go/helpers"
	"frank/src/go/helpers/log"
	"frank/src/go/models"
)

type Devices struct {
}

func ManageCommands(commandsFound []helpers.CommandFound) {
	if len(commandsFound) < 1 {
		return
	}

	for _, co := range commandsFound {
		log.Log.Debug("Handling command %s", co.CommandName)
		for _, a := range co.Actions {
			HandleActions(a, co.ExtraText)
		}
	}
}

func ManageReading(readingName string) (models.ReadingResponse, error) {
	reading, err := config.GetReading(readingName)
	if err != nil {
		log.Log.Warning(err.Error())
	}

	if reading.Plugin != "" {
		return ActivePlugins.ExecReading(reading)
	}

	log.Log.Warning("Could not handle Reading, Plugin attribute is missing", reading.Plugin)

	return models.ReadingResponse{}, errors.New("Could not find the plugin for the reading")
}

func ManageAction(actionName string, extraText map[string]string) {
	action, err := config.GetAction(actionName)
	if err != nil {
		log.Log.Warning(err.Error())
	}

	if action.Plugin != "" {
		ActivePlugins.ExecAction(action, extraText)
		return
	}

	log.Log.Warning("Could not handle Action, Plugin attribute is missing", action.Plugin)
}

func HandleActions(commandAction models.CommandAction, extraText map[string]string) {
	log.Log.Errorf("%+v AAAA %+v", commandAction, extraText)
	action, err := config.GetAction(commandAction.Action)
	if err != nil {
		log.Log.Warning(err.Error())
	}

	if action.Plugin != "" {
		ActivePlugins.ExecAction(action, extraText)
		return
	}

	log.Log.Warning("Could not handle Action, Plugin attribute is missing", action.Plugin)
}

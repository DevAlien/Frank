package devices

import (
	"frank/src/go/helpers/log"
	"frank/src/go/managers"
	"frank/src/go/models"
	"frank/src/go/services"
)

type Devices struct {
}

func ManageCommands(commandsFound []services.CommandFound) {
	if len(commandsFound) < 1 {
		return
	}

	for _, co := range commandsFound {
		log.Log.Debug("Handling command", co.CommandName)
		for _, a := range co.Actions {
			HandleActions(a, co.ExtraText)
		}
	}
}

func HandleActions(action models.CommandAction, extraText map[string]string) {
	if action.Plugin != "" {
		managers.ActivePlugins.ExecAction(action, extraText)
		return
	}

	log.Log.Warning("Could not handle Action, Plugin attribute is missing", action.Plugin)
}

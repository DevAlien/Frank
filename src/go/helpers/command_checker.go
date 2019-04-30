package helpers

import (
	"frank/src/go/models"
	"frank/src/go/config"
	"regexp"
)

type CommandFound struct {
	Text        string
	ExtraText   map[string]string
	CommandName string
	Actions     []models.CommandAction
}

var re = regexp.MustCompile(`{(.*?)}`)

func CheckCommands(text string) []CommandFound {
	commands := config.ParsedConfig.Commands
	var commandsFound []CommandFound
	for _, command := range commands {
		for _, textCommand := range command.Commands {
			s := re.ReplaceAllString(textCommand, `(?P<$1>.*)`)
			var r = regexp.MustCompile(s)
			match := r.FindStringSubmatch(text)
			if len(match) == 0 {
				continue
			}

			result := make(map[string]string)
			var tmpActions []models.CommandAction
			for i, name := range r.SubexpNames() {
				if i != 0 {
					result[name] = match[i]

					for _, a := range command.Actions {
						commandAction := a
						if len(a.MatchingAction) > 0 && a.Action == "" {

							if foundMatch, ok := a.MatchingAction[name]; ok {
								if val, ok := foundMatch.(map[string]interface{})[match[i]]; ok {
									commandAction.Action = val.(string)
									tmpActions = append(tmpActions, commandAction)
								}
							}
						}
					}
				}
			}
			if len(tmpActions) == 0 {
				tmpActions = command.Actions
			}
			commandsFound = append(commandsFound, CommandFound{
				Text:        text,
				ExtraText:   result,
				CommandName: command.Name,
				Actions:     tmpActions,
			})
		}
	}

	return commandsFound
}

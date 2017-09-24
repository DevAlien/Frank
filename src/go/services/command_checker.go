package services

import (
	"regexp"

	"frank/src/go/models"
)

type CommandFound struct {
	Text        string
	ExtraText   map[string]string
	CommandName string
	Actions     []models.CommandAction
}

var re = regexp.MustCompile(`{(.*?)}`)
var myExp = regexp.MustCompile(`asd {(?P<lol>.*)} var {(?P<mad>.*)}`)

func CheckCommands(text string, commands []models.Command) []CommandFound {
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
			for i, name := range r.SubexpNames() {
				if i != 0 {
					result[name] = match[i]
				}
			}

			commandsFound = append(commandsFound, CommandFound{
				Text:        text,
				ExtraText:   result,
				CommandName: command.Name,
				Actions:     command.Actions,
			})
		}
	}

	return commandsFound
}

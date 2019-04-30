package services

import (
	"strconv"
"time"
"strings"

	"frank/src/go/helpers"
	"frank/src/go/helpers/log"
	"frank/src/go/models"
	"frank/src/go/managers"
	
	"github.com/tucnak/telebot"
)

type TelegramBotService struct {
	Bot *telebot.Bot
	config *models.HTTP
}

var TelegramBot TelegramBotService

func NewTelegramBotService(config *models.Telegram) {
	if config.Disabled != true || config.APIKey == "" {
		return
	}

	log.Log.Info("Starting Telegram Bot")
	bot, err := telebot.NewBot(config.APIKey)
	if err != nil {
		log.Log.Critical(err.Error())
	}
	TelegramBot.Bot = bot
	TelegramBot.Bot.Messages = make(chan telebot.Message, 100)

	go TelegramBot.messages()

	go func() {
		TelegramBot.Bot.Start(1 * time.Second)
	}()
}

func (tbs *TelegramBotService) messages() {
	for message := range tbs.Bot.Messages {
		log.Log.Debugf("Received a message from %s with the text: %s %+v\n",
			message.Sender.Username, message.Text, message.Chat)
		commands := helpers.CheckCommands(strings.ToLower(message.Text))
		go managers.ManageCommands(commands)
		tbs.Bot.SendMessage(message.Sender,
			"Command Executed "+strconv.FormatBool(len(commands) > 0)+"!", nil)
	}
}

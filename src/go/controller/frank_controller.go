package controller

import (
	"fmt"
	"time"

	"frank/src/go/config"
	"frank/src/go/helpers"
	"frank/src/go/helpers/log"
	"frank/src/go/services"

	"frank/src/go/managers"
)

func task() {
	fmt.Println("I am runnning task.")
}

func taskWithParams(a int, b string) {
	fmt.Println(a, b)
}

const developerKey = "AIzaSyBEsKHzV5PkHUhvEOKjYfefv7_tkZ8EREs"
const (
	Stopped = 0
	Paused  = 1
	Running = 2
)

const timeout = 30 * time.Second

type FrankController struct {
}

func NewFrankController() (FrankController, error) {
	helpers.LoadDirs()

	frankController := FrankController{}
	log.InitLogger()

	log.Log.Warning("Frank is starting!!!")

	err := config.InitConfig(helpers.ConfigDir)
	if err != nil {
		log.Log.Critical(err.Error())
		return frankController, err
	}

	managers.NewPlugins()

	managers.LoadCrons()

	services.LoadServices()
	// if config.ParsedConfig.Ddns.Hostname != "" {
	// 	managers.LoadDdns(config.ParsedConfig.Ddns)
	// 	go managers.DdnsManager.SetIp()
	// }

	// if config.GetVoice().APIKey != "" && config.GetVoice().Disabled != true {
	// 	log.Log.Debug("Starting Voice Recognition Service")
	// 	voiceRecognition, _ := managers.NewVoiceRecognition(config.GetVoice())
	// 	frankController.VoiceRecognition = voiceRecognition
	// }

	// if config.GetWebSocket().Disabled != true {
	// 	frankController.SocketIoServer = managers.NewSocketIoServer(config.GetWebSocket())
	// }

	// if config.GetHTTP().Disabled != true {
	// 	managers.NewHttpServer(config.GetHTTP())
	// }

	// services.NewMqttServer(config.GetHTTP())
	return frankController, nil
}

// func (fc *FrankController) Start() {
// 	log.Log.Info(config.ParsedConfig.Name)
// 	fc.keywordCh = make(chan int)
// 	fc.voiceCh = make(chan int)
// 	fc.killCh = make(chan bool, 1)

// 	log.Log.Info("Starting Keyword Recognition")
// 	if config.Get("google_api_key") != "" && config.Is("voice_recognition") {
// 		go fc.StartKeywordRecognition()
// 	}

// 	// fc.SocketIoServer.Server.On("text", func(msg string) (bool, string) {
// 	// 	commands := helpers.CheckCommands(msg, config.ParsedConfig.Commands)
// 	// 	go managers.ManageCommands(commands)
// 	// 	return len(commands) > 0, "asd"
// 	// })

// 	managers.AddRoute("GET", "/command", func(c *gin.Context) {
// 		text := c.DefaultQuery("text", "")
// 		if text != "" {
// 			commands := helpers.CheckCommands(text, config.ParsedConfig.Commands)
// 			log.Log.Error("EEEE %+v", commands)
// 			go managers.ManageCommands(commands)
// 			a := "asd"
// 			c.JSON(200, a)
// 		}

// 	})

// 	fc.StartTelegramBot(config.Get("telegram_key"))

// 	var input string
// 	fmt.Scanln(&input)
// }

// func (fc *FrankController) StartTelegramBot(botKey string) {
// 	if botKey == "" {
// 		return
// 	}

// 	log.Log.Info("Starting Telegram Bot")
// 	bot, err := telebot.NewBot(botKey)
// 	if err != nil {
// 		log.Log.Critical(err.Error())
// 	}
// 	fc.Bot = bot
// 	fc.Bot.Messages = make(chan telebot.Message, 100)

// 	go fc.messages()

// 	go func() {
// 		fc.Bot.Start(1 * time.Second)
// 	}()
// }

// func (fc *FrankController) VoiceRecognitionToText(fileName string) {
// 	fc.SocketIoServer.Server.BroadcastTo("bot", "bot:analyzing", true)
// 	fmt.Printf("#goroutines: %d\n", runtime.NumGoroutine())
// 	log.Log.Debugf("[%s] Analyzing Audio", fileName)
// 	text := fc.VoiceRecognition.AnalyzeAudio(fileName)

// 	helpers.RemoveRecordFile(fileName)
// 	fc.SocketIoServer.Server.BroadcastTo("bot", "bot:text", text)
// 	log.Log.Debugf("[%s] Found Text: %s", fileName, text)
// 	commands := helpers.CheckCommands(text, config.ParsedConfig.Commands)
// 	go managers.ManageCommands(commands)
// 	fc.CheckDeactivation(fileName, text)
// 	fc.SocketIoServer.Server.BroadcastTo("bot", "bot:analyzing", false)
// 	go fc.StartVoiceRecognition()
// }

// func (fc *FrankController) CheckDeactivation(fileName string, text string) {
// 	log.Log.Debugf("[%s] Checking Deactivation Keywords", fileName)

// 	for _, sentence := range config.ParsedConfig.Deactivation {
// 		if sentence == text {
// 			log.Log.Infof("[%s] Deactivation Keyword Found", fileName)
// 			fc.StopVoiceRecognition()
// 			return
// 		}
// 	}
// 	fc.timer.Reset(timeout)
// }

// func (fc *FrankController) StopVoiceRecognition() {
// 	fc.SocketIoServer.Server.BroadcastTo("bot", "bot:listening", false)
// 	fc.SocketIoServer.Server.BroadcastTo("bot", "bot:sleep", true)
// 	log.Log.Debug("Stopping Voice Recognition And starting Keyword Recognition")
// 	fc.killCh <- true
// 	fc.voiceCh <- Stopped
// 	go fc.StartKeywordRecognition()
// }

// func (fc *FrankController) StartTimerStop() {
// 	if fc.timer != nil {
// 		fc.timer.Stop()
// 	}
// 	fc.timer = time.AfterFunc(timeout, fc.StopVoiceRecognition)
// }
// func (fc *FrankController) StartVoiceRecognition() {
// 	state := Running
// 	log.Log.Debug("Started timeout to deactivate voice recognition")
// 	fc.SocketIoServer.Server.BroadcastTo("bot", "bot:listening", true)
// 	fc.StartTimerStop()
// 	for {
// 		select {
// 		case state = <-fc.voiceCh:
// 			switch state {
// 			case Stopped:
// 				log.Log.Info("Stopped Voice Recognition")
// 				fc.keywordCh <- Running
// 				return
// 			case Running:
// 				log.Log.Info("Started Voice Recognition")
// 			case Paused:
// 				log.Log.Info("Paused Voice Recognition")
// 			}
// 		default:
// 			if state == Paused {
// 				time.Sleep(1 * time.Second)
// 				break
// 			}
// 			log.Log.Info("Listening Voice")
// 			fileName, _ := helpers.StartRecord(fc.killCh)
// 			if fileName == "" {
// 				break
// 			}
// 			fc.SocketIoServer.Server.BroadcastTo("bot", "bot:listening", false)
// 			go fc.VoiceRecognitionToText(fileName)
// 			return
// 		}
// 	}
// }
// func (fc *FrankController) StartKeywordRecognition() {
// 	state := Running
// 	for {
// 		select {
// 		case state = <-fc.keywordCh:
// 			switch state {
// 			case Stopped:
// 				log.Log.Info("Stopped Keyword Recognition")
// 				return
// 			case Running:
// 				log.Log.Info("Started Keyword Recognition")
// 			case Paused:
// 				log.Log.Info("Paused Keyword Recognition")
// 			}
// 		default:
// 			if state == Paused {
// 				time.Sleep(1 * time.Second)
// 				break
// 			}

// 			log.Log.Info("Listening Keyword")
// 			fileName, _ := helpers.StartRecord(fc.killCh)
// 			if fileName == "" {
// 				break
// 			}

// 			result := helpers.KeywordRecognition(fileName)
// 			helpers.RemoveRecordFile(fileName)
// 			if result == true {
// 				log.Log.Debug("Keyword matched")
// 				go fc.StartVoiceRecognition()
// 				return
// 			} else {
// 				log.Log.Debug("Keyword not matched")
// 				go fc.StartVoiceRecognition() //TODO REMOVE
// 				return
// 				// WHEN REMOVED THE ABOVE 2 LINES break
// 			}
// 		}
// 		time.Sleep(30)
// 	}
// }

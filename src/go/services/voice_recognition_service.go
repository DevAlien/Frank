package services

import (
	"fmt"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"runtime"

	"frank/src/go/helpers"
	"frank/src/go/managers"
	"frank/src/go/helpers/log"
	"frank/src/go/models"
	"frank/src/go/config"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/language/v1"
	"google.golang.org/api/speech/v1"
)

type VoiceRecognition struct {
	SpeechService   *speech.Service
	LanguageService *language.Service
	config *models.Voice
	timer            *time.Timer
	keywordCh        chan int
	voiceCh          chan int
	killCh           chan bool
}

const (
	Stopped = 0
	Paused  = 1
	Running = 2
)

const timeout = 30 * time.Second

var VR VoiceRecognition
func NewVoiceRecognition(config *models.Voice) {
	VR = VoiceRecognition{}
	VR.config = config

	speechService, err := speech.New(&http.Client{
		Transport: &transport.APIKey{Key: config.APIKey},
	})
	if err != nil {
		return
	}

	languageService, err := language.New(&http.Client{
		Transport: &transport.APIKey{Key: config.APIKey},
	})
	if err != nil {
		return
	}

	VR.SpeechService = speechService
	VR.LanguageService = languageService
	
	VR.Start()
}

func Encode(path string) (string, error) {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buff), nil
}

func (vr *VoiceRecognition) Start() {
	vr.keywordCh = make(chan int)
	vr.voiceCh = make(chan int)
	vr.killCh = make(chan bool, 1)
	log.Log.Info("Starting Keyword Recognition")
	if vr.config.APIKey != "" {
		log.Log.Info("Started Keyword Recognition")
		go vr.StartKeywordRecognition()
 	}

}


func (vr *VoiceRecognition) CheckDeactivation(fileName string, text string) {
	log.Log.Debugf("[%s] Checking Deactivation Keywords", fileName)

	for _, sentence := range config.ParsedConfig.Deactivation {
		if sentence == text {
			log.Log.Infof("[%s] Deactivation Keyword Found", fileName)
			vr.StopVoiceRecognition()
			return
		}
	}
	vr.timer.Reset(timeout)
}
func (vr *VoiceRecognition) VoiceRecognitionToText(fileName string) {
	SocketIo.Server.BroadcastTo("bot", "bot:analyzing", true)
	fmt.Printf("#goroutines: %d\n", runtime.NumGoroutine())
	log.Log.Debugf("[%s] Analyzing Audio", fileName)
	text := vr.AnalyzeAudio(fileName)

	helpers.RemoveRecordFile(fileName)
	SocketIo.Server.BroadcastTo("bot", "bot:text", text)
	log.Log.Debugf("[%s] Found Text: %s", fileName, text)
	commands := helpers.CheckCommands(text)
	go managers.ManageCommands(commands)
	vr.CheckDeactivation(fileName, text)
	SocketIo.Server.BroadcastTo("bot", "bot:analyzing", false)
	go vr.StartVoiceRecognition()
}

func (vr *VoiceRecognition) StopVoiceRecognition() {
	SocketIo.Server.BroadcastTo("bot", "bot:listening", false)
	SocketIo.Server.BroadcastTo("bot", "bot:sleep", true)
	log.Log.Debug("Stopping Voice Recognition And starting Keyword Recognition")
	vr.killCh <- true
	vr.voiceCh <- Stopped
	go vr.StartKeywordRecognition()
}

func (vr *VoiceRecognition) StartTimerStop() {
	if vr.timer != nil {
		vr.timer.Stop()
	}
	vr.timer = time.AfterFunc(timeout, vr.StopVoiceRecognition)
}
func (vr *VoiceRecognition) StartVoiceRecognition() {
	state := Running
	log.Log.Debug("Started timeout to deactivate voice recognition")
	SocketIo.Server.BroadcastTo("bot", "bot:listening", true)
	vr.StartTimerStop()
	for {
		select {
		case state = <-vr.voiceCh:
			switch state {
			case Stopped:
				log.Log.Info("Stopped Voice Recognition")
				vr.keywordCh <- Running
				return
			case Running:
				log.Log.Info("Started Voice Recognition")
			case Paused:
				log.Log.Info("Paused Voice Recognition")
			}
		default:
			if state == Paused {
				time.Sleep(1 * time.Second)
				break
			}
			log.Log.Info("Listening Voice")
			fileName, _ := helpers.StartRecord(vr.killCh)
			if fileName == "" {
				break
			}
			SocketIo.Server.BroadcastTo("bot", "bot:listening", false)
			go vr.VoiceRecognitionToText(fileName)
			return
		}
	}
}
func (vr *VoiceRecognition) StartKeywordRecognition() {
	state := Running
	for {
		select {
		case state = <-vr.keywordCh:
			switch state {
			case Stopped:
				log.Log.Info("Stopped Keyword Recognition")
				return
			case Running:
				log.Log.Info("Started Keyword Recognition")
			case Paused:
				log.Log.Info("Paused Keyword Recognition")
			}
		default:
			if state == Paused {
				time.Sleep(1 * time.Second)
				break
			}

			log.Log.Info("Listening Keyword")
			fileName, _ := helpers.StartRecord(vr.killCh)
			if fileName == "" {
				break
			}

			result := helpers.KeywordRecognition(fileName)
			helpers.RemoveRecordFile(fileName)
			if result == true {
				log.Log.Debug("Keyword matched")
				go vr.StartVoiceRecognition()
				return
			} else {
				log.Log.Debug("Keyword not matched")
				go vr.StartVoiceRecognition() //TODO REMOVE
				return
				// WHEN REMOVED THE ABOVE 2 LINES break
			}
		}
		time.Sleep(30)
	}
}


func (vr *VoiceRecognition) AnalyzeAudio(file string) string {
	start := time.Now()
	text, err := vr.sendAudioToGoogle(file)
	if err != nil {
		log.Log.Error("Error sending audio to google:", err)
	}

	// err = vr.parseText(text, file)
	// if err != nil {
	// 	log.Log.Error("Error parsing text:", err)
	// }

	elapsed := time.Since(start)
	log.Log.Debug("["+file+"] to get analysis", elapsed)

	return strings.ToLower(text)
}

func (vr *VoiceRecognition) sendAudioToGoogle(file string) (string, error) {
	var text string

	file64, err := Encode(helpers.GetRecordPath(file))
	if err != nil {
		return text, err
	}

	recognitionAudio := speech.RecognitionAudio{
		Content: file64,
	}

	recognitionConfig := speech.RecognitionConfig{
		LanguageCode:    "it-IT",
		Encoding:        "FLAC",
		SampleRateHertz: 16000,
	}

	recognizeRequest := speech.RecognizeRequest{
		Audio:  &recognitionAudio,
		Config: &recognitionConfig,
	}

	c := vr.SpeechService.Speech.Recognize(&recognizeRequest)
	response, err := c.Do()
	if err != nil {
		return text, err
	}

	for _, result := range response.Results {
		for _, alt := range result.Alternatives {
			text = alt.Transcript
			//fmt.Println("[" + file + "]", alt.Transcript, alt.Confidence)
		}
	}

	return text, nil
}

func (vr *VoiceRecognition) parseText(text string, file string) error {
	document := language.Document{
		Content:  text,
		Language: "it",
		Type:     "PLAIN_TEXT",
	}
	asr := language.AnalyzeSyntaxRequest{
		Document: &document,
	}
	call := vr.LanguageService.Documents.AnalyzeSyntax(&asr)

	response, err := call.Do()
	if err != nil {
		return err
	}

	for _, token := range response.Tokens {
		log.Log.Debug("["+file+"]", token.Text.Content, "->", token.Lemma, "=>", token.PartOfSpeech.Tag, token.PartOfSpeech.Form)
	}

	return nil
}

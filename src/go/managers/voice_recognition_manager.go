package managers

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"frank/src/go/helpers"
	"frank/src/go/helpers/log"
	"frank/src/go/models"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/language/v1"
	"google.golang.org/api/speech/v1"
)

type VoiceRecognition struct {
	SpeechService   *speech.Service
	LanguageService *language.Service
	config *models.Voice
}

func NewVoiceRecognition(config *models.Voice) (VoiceRecognition, error) {
	voiceRecognition := VoiceRecognition{}
	voiceRecognition.config = config

	speechService, err := speech.New(&http.Client{
		Transport: &transport.APIKey{Key: config.APIKey},
	})
	if err != nil {
		return voiceRecognition, err
	}

	languageService, err := language.New(&http.Client{
		Transport: &transport.APIKey{Key: config.APIKey},
	})
	if err != nil {
		return voiceRecognition, err
	}

	voiceRecognition.SpeechService = speechService
	voiceRecognition.LanguageService = languageService

	return voiceRecognition, nil
}

func Encode(path string) (string, error) {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buff), nil
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

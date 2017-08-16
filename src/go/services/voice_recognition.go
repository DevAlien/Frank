package services

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/base64"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/language/v1"
	"google.golang.org/api/speech/v1"
)

type VoiceRecognition struct {
	SpeechService *speech.Service
	LanguageService *language.Service
}

func NewVoiceRecognition(developerKey string) (VoiceRecognition, error) {
	voiceRecognition := VoiceRecognition{}

	speechService, err := speech.New(&http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	})
	if err != nil {
		return voiceRecognition, err
	}

	languageService, err := language.New(&http.Client{
		Transport: &transport.APIKey{Key: developerKey},
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
		log.Println("Error sending audio to google:", err)
	}

	err = vr.parseText(text, file)
	if err != nil {
		log.Println("Error parsing text:", err)
	}

	elapsed := time.Since(start)
  log.Println("[" + file + "] to get analysis", elapsed)

	return text
}

func (vr *VoiceRecognition) sendAudioToGoogle(file string) (string, error) {
	var text string

	file64, err := Encode("./" + file)
	if err != nil {
		return text, err
	}

	recognitionAudio := speech.RecognitionAudio{
		Content: file64,
	}

	recognitionConfig := speech.RecognitionConfig{
		LanguageCode: "it-IT",
		Encoding: "FLAC",
		SampleRateHertz: 16000,
	}

	recognizeRequest := speech.RecognizeRequest{
		Audio: &recognitionAudio,
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
			fmt.Println("[" + file + "]", alt.Transcript, alt.Confidence)
		}
	}

	return text, nil
}

func (vr *VoiceRecognition) parseText(text string, file string) error {
	document := language.Document{
		Content: text,
		Language: "it",
		Type: "PLAIN_TEXT",
	}
	asr := language.AnalyzeSyntaxRequest{
		Document: &document,
	}
	call := vr.LanguageService.Documents.AnalyzeSyntax(&asr)

	response, err := call.Do()
	if err != nil {
		return err
	}
	b, err := response.MarshalJSON()
	fmt.Println(string(b))
	for _, token := range response.Tokens {
		fmt.Println("[" + file + "]", token.Text.Content, "->", token.Lemma, "=>", token.PartOfSpeech.Tag, token.PartOfSpeech.Form)
	}

	return nil
}
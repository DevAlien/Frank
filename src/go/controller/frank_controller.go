package controller

import (
	"fmt"
	"runtime"
	"os"

	"frank/src/go/services"

)

const developerKey = "AIzaSyBEsKHzV5PkHUhvEOKjYfefv7_tkZ8EREs"
const (
    Stopped = 0
    Paused  = 1
    Running = 2
)

type FrankController struct {
	VoiceRecognition services.VoiceRecognition

	keywordCh chan int
	voiceCh chan int
}

func NewFrankController() (FrankController, error){
	frankController := FrankController{}

	voiceRecognition, _ := services.NewVoiceRecognition(developerKey)
	frankController.VoiceRecognition = voiceRecognition

	return frankController, nil
}

func (fc *FrankController) Start() {
	fc.keywordCh = make(chan int)
	fc.voiceCh = make(chan int)
	go fc.StartKeywordRecognition()

	var input string
  fmt.Scanln(&input)
}

func (fc *FrankController) VoiceRecognitionToText(fileName string) {
	text := fc.VoiceRecognition.AnalyzeAudio(fileName)
	_ = os.Remove(fileName)

	fmt.Println(text)
}

func (fc *FrankController) StartVoiceRecognition() {
	state := Running
	for {
      select {
      case state = <- fc.voiceCh:
        switch state {
				case Stopped:
						fmt.Printf("Worker: Stopped\n")
						fc.keywordCh <- Running
						return
				case Running:
						fmt.Printf("Worker: Running\n")
				case Paused:
						fmt.Printf("Worker: Paused\n")
				}
			default:
				runtime.Gosched()
				if state == Paused {
						break
				}
				fileName, _ := services.StartRecord()
				go fc.VoiceRecognitionToText(fileName)
      }
    }
}
func (fc *FrankController) StartKeywordRecognition() {
	state := Running
	fmt.Println("asd")
	for {
      select {
      case state = <- fc.keywordCh:
        switch state {
				case Stopped:
						fmt.Printf("Worker: Stopped\n")
						return
				case Running:
						fmt.Printf("Worker: Running\n")
				case Paused:
						fmt.Printf("Worker: Paused\n")
				}
			default:
				runtime.Gosched()
				if state == Paused {
						break
				}
				fmt.Println("looking for Frank")
				fileName, _ := services.StartRecord()
				result := services.KeywordRecognition(fileName)
				_ = os.Remove(fileName)
				if result == true {
					fmt.Println("it is true")
					go fc.StartVoiceRecognition()
					state = Paused
					break
				} else {
					fmt.Println("it is false")
					go fc.StartVoiceRecognition()
					state = Paused
					break
				}
      }
    }
}
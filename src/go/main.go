package main

import (
	"fmt"
	"time"
	"log"
	"net/http"
	"os/exec"
	"bytes"
	"io/ioutil"
	"encoding/base64"
	"unsafe"
	
	"frank/src/go/controller"

	"github.com/satori/go.uuid"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/language/v1"
	"google.golang.org/api/speech/v1"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"

	"github.com/Kitt-AI/snowboy/swig/Go"
)
const developerKey = "AIzaSyBEsKHzV5PkHUhvEOKjYfefv7_tkZ8EREs"

func Encode(path string) (string, error) {
	buff, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buff), nil
}

func main() {
	fc, _ := controller.NewFrankController()
	fc.Start()
}

func mainQQ() {
				i := 0
				gobot.Every(1*time.Second, func() {
					if i == 1 {
						i = 0
					} else {
						i = 1
					}
					firmataAdaptor := firmata.NewTCPAdaptor("192.168.1.7:3030")
					_ = firmataAdaptor.Connect()
					fmt.Println("change", i)
					firmataAdaptor.DigitalWrite("12", byte(i))
				})
        // led := gpio.NewLedDriver(firmataAdaptor, "12")

        // work := func() {
        //         gobot.Every(5*time.Second, func() {
				// 					fmt.Println("toggling")
        //                 led.Toggle()
        //         })
        // }

        // robot := gobot.NewRobot("bot",
        //         []gobot.Connection{firmataAdaptor},
        //         []gobot.Device{led},
        //         work,
        // )

				// robot.Start()
				var input string
	fmt.Scanln(&input)
}

// func main2d() {
// 	voiceRecognition, _ := services.NewVoiceRecognition(developerKey)
// 	fileName, _ := services.StartRecord()
// 	text := voiceRecognition.AnalyzeAudio(fileName)
// 	_ = os.Remove(fileName)
// 	fmt.Println("[" + fileName+ "] Text: ", text)
// }

func mainq() {
	fileName := fmt.Sprintf("%s.flac", uuid.NewV4())
	fmt.Println("[" + fileName+ "] listening...")
	
	cmd := exec.Command("rec", "-r", "16000", "-c", "1", fileName, "silence", "-l", "1", "0.5", "0.1%", "1", "1.0", "0.1%")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(out.String())
		log.Fatal(err)
	}
	
	
	detector := snowboydetect.NewSnowboyDetect("./common.res", "./Frank.pmdl")
	detector.SetSensitivity("0.5")
	// detector.SetAudioGain(1)
	defer snowboydetect.DeleteSnowboyDetect(detector)

	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	ptr := snowboydetect.SwigcptrInt16_t(unsafe.Pointer(&dat[0]))
	res := detector.RunDetection(ptr, len(dat) / 2 /* len of int16  */)
	if res == -2 {
		fmt.Println("Snowboy detected silence")
	} else if res == -1 {
		fmt.Println("Snowboy detection returned error")
	} else if res == 0 {
		fmt.Println("Snowboy detected nothing")
	} else {
		fmt.Println("Snowboy detected keyword ", res)
	}
}
func main22() {
	c1 := make(chan string)

	firmataAdaptor := firmata.NewAdaptor("/dev/tty.usbmodem1411")
	led := gpio.NewLedDriver(firmataAdaptor, "13")
	led2 := gpio.NewLedDriver(firmataAdaptor, "12")

	work := func() {
					// gobot.Every(1*time.Second, func() {
					// 				led.Toggle()
					// })
	}

	robot := gobot.NewRobot("bot",
					[]gobot.Connection{firmataAdaptor},
					[]gobot.Device{led, led2},
					work,
	)

	
	
  go func() {
    for {
			fileName := fmt.Sprintf("%s.flac", uuid.NewV4())
			fmt.Println("[" + fileName+ "] listening...")
			
			cmd := exec.Command("rec", "-r", "16000", "-c", "1", fileName, "silence", "-l", "1", "0.5", "0.1%", "1", "1.0", "0.1%")
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			if err != nil {
				fmt.Println(out.String())
				log.Fatal(err)
			}
      c1 <- fileName
    }
  }()

  go func() {
    for {
      select {
      case msg1 := <- c1:
        analyzeAudio(msg1, led, led2)
      }
    }
  }()
	robot.Start()
	var input string
  fmt.Scanln(&input)
}

func analyzeAudio(file string, led *gpio.LedDriver, led2 *gpio.LedDriver) {
	start := time.Now()
	text, err := sendAudioToGoogle(file)
	if err != nil {
		log.Println("Error sending audio to google: %v", err)
	}

	if text == "accendi la luce verde" {
		led.On()
	} else if text == "spegni la luce verde" {
		led.Off()
	}

	if text == "accendi la luce blu" {
		fmt.Println("on blu")
		led2.On()
	} else if text == "spegni la luce blu" {
		fmt.Println("off blu")
		led2.Off()
	}

	if text == "accendi le luci" {
		led2.On()
		led.On()
	} else if text == "spegni le luci" {
		led2.Off()
		led.Off()
	}


	err = parseText(text, file)
	if err != nil {
		log.Println("Error parsing text: %v", err)
	}

	elapsed := time.Since(start)
  log.Println("[" + file + "] to get analysis %s", elapsed)

}

func sendAudioToGoogle(file string) (string, error) {
	var text string
	speechService, err := speech.New(&http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	})
	if err != nil {
		return text, err
	}

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

	c := speechService.Speech.Recognize(&recognizeRequest)
	response1, err := c.Do()
	if err != nil {
		return text, err
	}

	for _, result := range response1.Results {
		for _, alt := range result.Alternatives {
			text = alt.Transcript
			fmt.Println("[" + file + "]", alt.Transcript, alt.Confidence)
		}
	}

	return text, nil
}

func parseText(text string, file string) error {
	languageService, err := language.New(&http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	})
	if err != nil {
		return err
	}

	document := language.Document{
		Content: text,
		Language: "it",
		Type: "PLAIN_TEXT",
	}
	asr := language.AnalyzeSyntaxRequest{
		Document: &document,
	}
	call := languageService.Documents.AnalyzeSyntax(&asr)

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

func main2() {
	languageService, err := language.New(&http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	})
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	speechService, err := speech.New(&http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	})
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	cmd := exec.Command("rec", "-r", "16000", "-c", "1", "record2.flac", "silence", "-l", "1", "0.5", "0.1%", "1", "1.0", "0.1%")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		fmt.Println(out.String())
		log.Fatal(err)
	}
	fmt.Println(out)
	start := time.Now()

	
	file64, err := Encode("/Users/goncalo/go/src/frank/src/go/record2.flac")
	if err != nil {
		log.Fatalf("file conversion: %v", err)
	}
	elapsed := time.Since(start)
  log.Printf("to base 64 %s", elapsed)
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

	c := speechService.Speech.Recognize(&recognizeRequest)
	response1, err := c.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}
	elapsed = time.Since(start)
	var text string

  log.Printf("to retrieve text %s", elapsed)
	for _, result := range response1.Results {
		for _, alt := range result.Alternatives {
			text = alt.Transcript
			fmt.Println(alt.Transcript, alt.Confidence)
		}
	}

	document := language.Document{
		Content: text,
		Language: "it",
		Type: "PLAIN_TEXT",
	}
	asr := language.AnalyzeSyntaxRequest{
		Document: &document,
	}
	call := languageService.Documents.AnalyzeSyntax(&asr)

	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}
	elapsed = time.Since(start)
  log.Printf("to get analysis %s", elapsed)

	for _, token := range response.Tokens {
		fmt.Println(token.Text.Content, "->", token.Lemma, "=>", token.PartOfSpeech.Tag, token.PartOfSpeech.Form)
	}
	fmt.Printf("%+v\n", response.Tokens)
	fmt.Println(developerKey)
}
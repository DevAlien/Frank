package services

import (
	"fmt"
	"unsafe"
	"io/ioutil"

	"github.com/Kitt-AI/snowboy/swig/Go"
)

const restFile = "./common.res"
const pmflFile = "./Frank.pmdl"

func KeywordRecognition(fileName string) bool{

	detector := snowboydetect.NewSnowboyDetect(restFile, pmflFile)
	detector.SetSensitivity("0.5")
	// detector.SetAudioGain(1)
	defer snowboydetect.DeleteSnowboyDetect(detector)

	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return false
	}

	ptr := snowboydetect.SwigcptrInt16_t(unsafe.Pointer(&dat[0]))
	res := detector.RunDetection(ptr, len(dat) / 2 /* len of int16  */)
	if res == -2 {
		fmt.Println("Snowboy detected silence")
		return false
	} else if res == -1 {
		fmt.Println("Snowboy detection returned error")
		return false
	} else if res == 0 {
		fmt.Println("Snowboy detected nothing")
		return false
	} else {
		fmt.Println("Snowboy detected keyword ", res)
		return true
	}
}
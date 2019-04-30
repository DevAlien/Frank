package helpers

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"frank/src/go/config"
	"frank/src/go/helpers/log"

	"github.com/satori/go.uuid"
)

const defaultCommand = "-r 16000 -c 1 %s silence " // dirFile, "silence"
const silenceDefault = "-l 1 0.2 2% 1 0.2 2%"

func StartRecord(killChannel chan bool) (string, error) {
	fileName := fmt.Sprintf("%s.flac", uuid.NewV4())
	log.Log.Info("[" + fileName + "] listening...")

	silenceParams := config.Get("record_silence_params")
	if silenceParams == "" {
		log.Log.Debug("using default silence")
		silenceParams = silenceDefault
	}

	dirFile := GetRecordPath(fileName)
	parts := strings.Fields(fmt.Sprintf(defaultCommand, dirFile) + silenceParams)

	cmd := exec.Command("rec", parts...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Start()
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-killChannel:
		if err := cmd.Process.Kill(); err != nil {
			return fileName, err
		}
		_ = os.Remove(dirFile)

		return "", nil
	case err := <-done:
		if err != nil {
			log.Log.Critical("["+fileName+"]", err)
			return fileName, err
		} else {
			log.Log.Info("[" + fileName + "] received voice")
			return fileName, nil
		}
	}
	return fileName, nil
}

package services

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"frank/src/go/helpers"
	"frank/src/go/helpers/log"

	"github.com/satori/go.uuid"
)

func StartRecord(killChannel chan bool) (string, error) {
	fileName := fmt.Sprintf("%s.flac", uuid.NewV4())
	log.Log.Info("[" + fileName + "] listening...")
	dirFile := helpers.GetRecordPath(fileName)
	cmd := exec.Command("rec", "-r", "16000", "-c", "1", dirFile, "silence", "-l", "1", "0.5", "0.1%", "1", "1.0", "0.1%")
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

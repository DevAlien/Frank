package services

import (
	"fmt"
	"os/exec"
	"bytes"

	"frank/src/go/helpers/log"

	"github.com/satori/go.uuid"
)

func StartRecord(killChannel chan bool) (string, error){
	fileName := fmt.Sprintf("%s.flac", uuid.NewV4())
	log.Log.Info("[" + fileName+ "] listening...")
	
	cmd := exec.Command("rec", "-r", "16000", "-c", "1", fileName, "silence", "-l", "1", "0.5", "0.1%", "1", "1.0", "0.1%")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Start()
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	log.Log.Critical("Select")
	select {
	case state := <-killChannel:
		log.Log.Critical("IT IS KILLED", state)
		if err := cmd.Process.Kill(); err != nil {
				return fileName, err
		}
		return "", nil
	case err := <-done:
		log.Log.Critical("IT IS DONE", err)
		if err != nil {
				return fileName, err
		} else {
				log.Log.Info("[" + fileName+ "] received voice")
				return fileName, nil
		}
	}
	log.Log.Critical("After Select")
	return fileName, nil
}
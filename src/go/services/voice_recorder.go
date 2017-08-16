package services

import (
	"fmt"
	"os/exec"
	"bytes"

	"github.com/satori/go.uuid"
)

func StartRecord() (string, error){
	fileName := fmt.Sprintf("%s.flac", uuid.NewV4())
	fmt.Println("[" + fileName+ "] listening...")
	
	cmd := exec.Command("rec", "-r", "16000", "-c", "1", fileName, "silence", "-l", "1", "0.5", "0.1%", "1", "1.0", "0.1%")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fileName, nil
	}

	return fileName, nil
}
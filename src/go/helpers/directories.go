package helpers

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"frank/src/go/templates"
)

var BaseDir string
var TmpDir string
var ConfigDir string

const privateDirName = ".frank"

func LoadDirs() {
	setDirs()
	if _, err := os.Stat(TmpDir); os.IsNotExist(err) {
		createFiles()
	}

	removeContents(TmpDir)
}

func GetRecordPath(fileName string) string {
	return filepath.Join(TmpDir, fileName)
}

func RemoveRecordFile(fileName string) {
	os.Remove(GetRecordPath(fileName))
}
func setDirs() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("asd", err)
	}
	BaseDir = filepath.Join(usr.HomeDir, privateDirName)

	TmpDir = filepath.Join(BaseDir, "tmp")
	ConfigDir = filepath.Join(BaseDir, "config.json")
}

func createFiles() {
	createAndSetDirs()
	configTemplate := []byte(templates.ConfigTemplate)
	err := ioutil.WriteFile(ConfigDir, configTemplate, 0644)
	if err != nil {
		log.Fatal("asd1", err)
	}
}

func createAndSetDirs() {
	os.MkdirAll(TmpDir, os.ModePerm)
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

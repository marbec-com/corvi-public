package controllers

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"marb.ec/corvi-backend/models"
	"os"
	"sync"
)

const (
	settingsFileName string = "settings.yml"
)

var SettingsControllerSingleton *SettingsController

type SettingsController struct {
	settings *models.Settings
	sync.Mutex
}

func SettingsControllerInstance() *SettingsController {
	if SettingsControllerSingleton == nil {
		SettingsControllerSingleton = NewSettingsController()
	}
	return SettingsControllerSingleton
}

func NewSettingsController() *SettingsController {
	settingsData := models.NewSettings()
	settingsFile, err := os.Open(settingsFileName)
	defer settingsFile.Close()

	if err != nil && os.IsNotExist(err) {
		// Settings file does not exist, create new one
		err = createNewSettingsFile(settingsData)
		if err != nil {
			log.Fatalf("Error while creating a new settings file: %v.\n", err)
		}

	} else if err == nil {
		settingsData, err = loadFromSettingsFile(settingsFile)
		// TODO(mjb): Reset settings file in case of error
		if err != nil {
			log.Fatalf("Error while loading settings file: %v.\n", err)
		}
	} else {
		// Unknown error while opening file
		log.Fatalf("Error while opening settings file: %v.\n", err)
	}

	return &SettingsController{
		settings: settingsData,
	}
}

func (c *SettingsController) Update() error {
	c.Lock()
	defer c.Unlock()

	settingsFile, err := os.OpenFile(settingsFileName, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	return saveToFile(c.settings, settingsFile)
}

func (c *SettingsController) Get() *models.Settings {
	return c.settings
}

func saveToFile(settings *models.Settings, file *os.File) error {
	data, err := yaml.Marshal(settings)
	if err != nil {
		return err
	}

	n, err := file.Write(data)
	if err != nil {
		return err
	}
	if n < len(data) {
		return io.ErrShortWrite
	}

	return nil
}

func createNewSettingsFile(settings *models.Settings) error {
	settingsFile, err := os.Create(settingsFileName)
	if err != nil {
		return err
	}

	return saveToFile(settings, settingsFile)
}

func loadFromSettingsFile(file *os.File) (*models.Settings, error) {
	// Settings file exists, parse into settingsData
	buffer := bytes.NewBuffer(nil)

	_, err := io.Copy(buffer, file)
	if err != nil {
		return nil, err
	}

	settings := models.NewSettings()
	err = yaml.Unmarshal(buffer.Bytes(), &settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

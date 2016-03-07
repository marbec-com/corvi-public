package controllers

import (
	"log"
	"os"
	"path"
)

func InitControllerSingletons(db *DBController) {
	settingsFileName := GenerateUserDataPath(settingsFile)
	SettingsControllerSingleton = NewSettingsController(settingsFileName)

	BoxControllerSingleton = NewBoxController(db)

	c, err := NewCategoryController(db)
	if err != nil {
		log.Fatal(err)
	}
	CategoryControllerSingleton = c

	QuestionControllerSingleton = NewQuestionController(db)
	StatsControllerSingleton = NewStatsController(db)
}

func GenerateUserDataPath(fileName string) string {
	userPath := os.Getenv("USER_DATA")
	if userPath != "" {
		fileName = path.Join(userPath, fileName)
	}
	return fileName
}

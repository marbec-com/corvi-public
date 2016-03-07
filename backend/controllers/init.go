package controllers

import (
	"log"
	"os"
	"path"
)

func InitControllerSingletons(db *DBController) {
	settingsFileName := GenerateUserDataPath(settingsFile)
	SettingsControllerSingleton = NewSettingsController(settingsFileName)

	b, err := NewBoxController(db)
	if err != nil {
		log.Fatal(err)
	}
	BoxControllerSingleton = b

	c, err := NewCategoryController(db)
	if err != nil {
		log.Fatal(err)
	}
	CategoryControllerSingleton = c

	q, err := NewQuestionController(db)
	if err != nil {
		log.Fatal(err)
	}
	QuestionControllerSingleton = q

	StatsControllerSingleton = NewStatsController(db)
}

func GenerateUserDataPath(fileName string) string {
	userPath := os.Getenv("USER_DATA")
	if userPath != "" {
		fileName = path.Join(userPath, fileName)
	}
	return fileName
}

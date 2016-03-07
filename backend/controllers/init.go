package controllers

import (
	"os"
	"path"
)

func InitControllerSingletons(db *DBController) {
	userPath := os.Getenv("USER_DATA")
	settingsFileName := settingsFile
	if userPath != "" {
		settingsFileName = path.Join(userPath, settingsFile)
	}
	SettingsControllerSingleton = NewSettingsController(settingsFileName)

	BoxControllerSingleton = NewBoxController(db)
	CategoryControllerSingleton = NewCategoryController(db)
	QuestionControllerSingleton = NewQuestionController(db)
	StatsControllerSingleton = NewStatsController(db)
}

package models

type Settings struct {
	MaxDailyQuestionsPerBox  uint
	RelearnUntilAccomplished bool
}

func NewSettings() *Settings {
	// Default parameters
	return &Settings{
		MaxDailyQuestionsPerBox:  20,
		RelearnUntilAccomplished: false,
	}
}

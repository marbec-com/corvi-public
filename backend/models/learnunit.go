package models

import (
	"time"
)

type LearnUnit struct {
	QuestionID  uint
	BoxID       uint
	CreatedAt   time.Time
	Correct     bool
	PrevCorrect bool
}

func NewLearnUnit() *LearnUnit {
	return &LearnUnit{
		CreatedAt: time.Now(),
	}
}

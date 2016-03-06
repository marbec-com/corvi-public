package models

import (
	"time"
)

type LearnUnit struct {
	BoxID       uint
	QuestionID  uint
	CreatedAt   time.Time
	Correct     bool
	PrevCorrect bool
}

func NewLearnUnit(qID, boxID uint, c, prev bool) *LearnUnit {
	return &LearnUnit{
		BoxID:       boxID,
		QuestionID:  qID,
		CreatedAt:   time.Now(),
		Correct:     c,
		PrevCorrect: prev,
	}
}

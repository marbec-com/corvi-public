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

func NewLearnUnit(qID, boxID uint, c, prev bool) *LearnUnit {
	return &LearnUnit{
		QuestionID:  qID,
		BoxID:       boxID,
		CreatedAt:   time.Now(),
		Correct:     c,
		PrevCorrect: prev,
	}
}

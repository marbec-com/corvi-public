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

func (u *LearnUnit) Equal(a *LearnUnit) bool {
	return u.QuestionID == a.QuestionID && u.BoxID == a.BoxID && u.CreatedAt.Equal(a.CreatedAt) && u.Correct == a.Correct && u.PrevCorrect == a.PrevCorrect
}

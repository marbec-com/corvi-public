package models

import (
	"time"
)

type LearnUnit struct {
	BoxID       uint
	Box         *Box `json:"-"`
	QuestionID  uint
	Question    *Question `json:"-"`
	Time        time.Time
	Correct     bool
	PrevCorrect bool
}

func NewLearnUnit(q *Question, c, prev bool) *LearnUnit {
	return &LearnUnit{
		BoxID:       q.BoxID,
		Box:         q.Box,
		QuestionID:  q.ID,
		Question:    q,
		Time:        time.Now(),
		Correct:     c,
		PrevCorrect: prev,
	}
}

package models

import (
	"time"
)

type LearnUnit struct {
	Box      *Box
	Question *Question
	Time     time.Time
	Correct  bool
}

func NewLearnUnit(q *Question, c bool) *LearnUnit {
	return &LearnUnit{
		Box:      q.Box,
		Question: q,
		Time:     time.Now(),
		Correct:  c,
	}
}

package models

import (
	"time"
)

type Box struct {
	ID               uint
	Name             string
	Description      string
	CategoryID       uint
	QuestionHeap     *QuestionHeap `json:"-"`
	QuestionsTotal   uint
	QuestionsToLearn uint
	QuestionsLearned uint
	CreatedAt        time.Time
}

/*
 Published bool
*/

/*
PublishedBox
Box
Author
LastChange
Description



*/

func NewBox() *Box {
	return &Box{
		QuestionHeap: NewQuestionHeap(),
		CreatedAt:    time.Now(),
	}
}

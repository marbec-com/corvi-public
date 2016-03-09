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

func (b *Box) Equal(a *Box) bool {
	return b.ID == a.ID && b.Name == a.Name && b.Description == a.Description && b.CategoryID == a.CategoryID && b.QuestionHeap.Equal(a.QuestionHeap) && b.QuestionsTotal == a.QuestionsTotal && b.QuestionsToLearn == a.QuestionsToLearn && b.QuestionsLearned == a.QuestionsLearned && b.CreatedAt.Equal(a.CreatedAt)
}

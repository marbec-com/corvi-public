package models

type Box struct {
	ID               uint
	Name             string
	CategoryID       uint
	Category         *Category     `json:"-"`
	QuestionHeap     *QuestionHeap `json:"-"`
	QuestionsTotal   uint
	QuestionsToLearn uint
	QuestionsLearned uint
}

/*
 Published bool
 Description string
*/

/*
PublishedBox
Box
Author
LastChange
Description



*/

func NewBox(name string, cat *Category) *Box {
	return &Box{
		Name:         name,
		CategoryID:   cat.ID,
		Category:     cat,
		QuestionHeap: NewQuestionHeap(),
	}
}

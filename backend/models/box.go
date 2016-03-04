package models

type Box struct {
	ID               uint
	Name             string
	Description      string
	CategoryID       uint
	QuestionHeap     *QuestionHeap `json:"-"`
	QuestionsTotal   uint
	QuestionsToLearn uint
	QuestionsLearned uint
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
	}
}

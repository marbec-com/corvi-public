package models

type Box struct {
	ID               uint
	Name             string
	Category         *Category
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
		Category:     cat,
		QuestionHeap: NewQuestionHeap(),
	}
}

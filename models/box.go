package models

type Box struct {
	ID               uint
	Name             string
	Category         *Category
	QuestionHeap     *QuestionHeap `json:"-"`
	QuestionsTotal   int
	QuestionsToLearn int
	QuestionsLearned int
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

const (
	maxToLearn uint = 20
)

func NewBox(name string, cat *Category) *Box {
	return &Box{
		Name:         name,
		Category:     cat,
		QuestionHeap: NewQuestionHeap(),
	}
}

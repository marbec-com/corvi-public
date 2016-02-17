package models

type Box struct {
	ID                 uint
	Name               string
	Category           *Category
	QuestionsToLearn   *QuestionHeap
	QuestionsTotal     int
	QuestionsUnlearned int
	QuestionsLearned   int
}

const (
	maxToLearn uint = 20
)

func NewBox(name string, cat *Category) *Box {
	return &Box{
		Name:             name,
		Category:         cat,
		QuestionsToLearn: NewQuestionHeap(),
	}
}

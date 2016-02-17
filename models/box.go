package models

import (
	"errors"
	"time"
)

type Box struct {
	ID               uint
	Name             string
	Category         *Category
	questionsToLearn *QuestionHeap
}

const (
	maxToLearn uint = 20
)

func (b *Box) Questions() []*Question {

	q := make([]*Question, 0)

	// SELECT * FROM Question WHERE BoxID = b.ID

	return q
}

func (b *Box) QuestionsLearnedNumber() uint {

	// SELECT COUNT(*) FROM Questions WHERE BoxID = b.ID AND CorrectlyAnswered > 0

	return 0
}

func (b *Box) QuestionsUnlearnedNumber() uint {

	// SELECT COUNT(*) FROM Questions WHERE BoxID = b.ID AND CorrectlyAnswered = 0

	return 0
}

func (b *Box) QuestionsTotal() uint {

	// SELECT COUNT(*) FROM Questions WHERE BoxID = b.ID

	return 0

}

func (b *Box) GetQuestionToLearn() (*Question, error) {
	if b.questionsToLearn.Length() == 0 {
		b.loadQuestionsToLearn()
	}

	question := b.questionsToLearn.Min()

	if question.Next.After(time.Now()) { // TODO: Update to give questions for whole day
		return nil, errors.New("No new Questions for today.")
	}

	return question, nil

}

func (b *Box) loadQuestionsToLearn() {

	// SELECT * FROM Questions WHERE BoxID = b.ID ORDER BY Next DESC

	//b.questionsToLearn.Add()

}

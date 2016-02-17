package controllers

import (
	"marb.ec/corvi-backend/models"
)

type QuestionController struct {
}

func (c *QuestionController) LoadQuestions() ([]*models.Question, error) {

	return nil, nil
}

func (c *QuestionController) LoadQuestion(id uint) (*models.Question, error) {

	return nil, nil
}

func (c *QuestionController) GiveCorrectAnswer(id uint) error {

	// Update Numbers in box

	question, err := c.LoadQuestion(id)
	if err != nil {
		return err
	}
	question.CorrectlyAnswered++
	question.CalculateNext()

	return nil
}

func (c *QuestionController) GiveWrongAnswer(id uint) error {

	question, err := c.LoadQuestion(id)
	if err != nil {
		return err
	}
	question.CorrectlyAnswered = 0
	question.CalculateNext()

	// Update Heap and Numbers in box

	return nil
}

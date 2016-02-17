package controllers

import (
	"marb.ec/corvi-backend/models"
)

var QuestionControllerSingleton *QuestionController

type QuestionController struct {
}

func QuestionControllerInstance() *QuestionController {
	if QuestionControllerSingleton == nil {
		QuestionControllerSingleton = NewQuestionController()
	}
	return QuestionControllerSingleton
}

func NewQuestionController() *QuestionController {
	return &QuestionController{}
}

func (c *QuestionController) LoadQuestions() ([]*models.Question, error) {

	return nil, nil
}

func (c *QuestionController) LoadQuestion(id uint) (*models.Question, error) {
	// Get box from cache, otherwise load

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

	// Save Question
	err = c.UpdateQuestion(question)
	if err != nil {
		return err
	}

	return nil
}

func (c *QuestionController) GiveWrongAnswer(id uint) error {

	question, err := c.LoadQuestion(id)
	if err != nil {
		return err
	}
	question.CorrectlyAnswered = 0
	question.CalculateNext()

	// Save Question
	err = c.UpdateQuestion(question)
	if err != nil {
		return err
	}
	// Update Heap and Numbers in box

	return nil
}

func (c *QuestionController) UpdateQuestion(q *models.Question) error {
	return nil
}

func (c *QuestionController) InsertQuestion(q *models.Question) error {
	// Update Box
	return nil
}

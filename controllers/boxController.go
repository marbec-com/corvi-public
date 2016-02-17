package controllers

import (
	"errors"
	"marb.ec/corvi-backend/models"
	"time"
)

// TODO: Monitor Memory Consumption
var BoxCache map[uint]*models.Box

type BoxController struct {
}

func (c *BoxController) LoadBoxes() ([]*models.Box, error) {

	// Load all form SQL

	// Save in BaxCache

	return nil, nil
}

func (c *BoxController) LoadBox(id uint) (*models.Box, error) {
	box, ok := BoxCache[id]
	if ok {
		return box, nil
	}

	// Load from SQL

	return nil, nil
}

func (c *BoxController) LoadQuestions(id uint) (*[]models.Question, error) {

	return nil, nil
}

func (c *BoxController) GetQuestionToLearn(id uint) (*models.Question, error) {
	box, err := c.LoadBox(id)
	if err != nil {
		return nil, err
	}

	if box.QuestionsToLearn.Length() == 0 {
		c.loadQuestionsToLearn(box)
	}

	question := box.QuestionsToLearn.Min()

	if question.Next.After(time.Now()) { // TODO: Update to give questions for whole day
		return nil, errors.New("No question to learn for this box.")
	}

	return question, nil
}

func (c *BoxController) loadQuestionsToLearn(b *models.Box) {

	// SELECT * FROM Questions WHERE BoxID = b.ID ORDER BY Next DESC

	//b.questionsToLearn.Add()

}

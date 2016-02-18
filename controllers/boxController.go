package controllers

import (
	"errors"
	"marb.ec/corvi-backend/models"
)

var BoxControllerSingleton *BoxController

type BoxController struct {
	BoxCache map[uint]*models.Box // TODO: Monitor Memory Consumption
}

func BoxControllerInstance() *BoxController {
	if BoxControllerSingleton == nil {
		BoxControllerSingleton = NewBoxController()
	}
	return BoxControllerSingleton
}

func NewBoxController() *BoxController {
	return &BoxController{}
}

func (c *BoxController) LoadBoxes() ([]*models.Box, error) {

	// Load all form SQL

	// Save in BaxCache

	return []*models.Box{}, nil
}

func (c *BoxController) LoadBox(id uint) (*models.Box, error) {
	box, ok := c.BoxCache[id]
	if ok {
		return box, nil
	}

	// Load Box SQL

	// Store in Cache
	//BoxCache[box.ID] = box

	return &models.Box{}, nil
}

func (c *BoxController) refreshBox(box *models.Box) error {

	// Subqueries for numbers

	return nil
}

func (c *BoxController) LoadQuestions(id uint) ([]*models.Question, error) {

	return []*models.Question{}, nil
}

func (c *BoxController) GetQuestionToLearn(id uint) (*models.Question, error) {
	box, err := c.LoadBox(id)
	if err != nil {
		return nil, err
	}

	// Load more questions if heap is nil.
	if box.QuestionsToLearn.Length() == 0 {
		c.loadQuestionsToLearn(box)
	}

	// Get next question in heap. Might be nil if no questions are remaining for this day
	question := box.QuestionsToLearn.Min()

	if question == nil {
		return nil, errors.New("No question to learn for this box.")
	}

	return question, nil
}

func (c *BoxController) loadQuestionsToLearn(b *models.Box) {

	// SELECT * FROM Questions WHERE BoxID = b.ID ORDER BY Next DESC LIMIT 20

	// Skip if question.Next.After(time.Now())
	// TODO: Update to give questions for whole day

	//b.questionsToLearn.Add()

}

func (c *BoxController) UpdateBox(b *models.Box) error {
	c.BoxCache[b.ID] = b
	return nil
}

func (c *BoxController) Insert(b *models.Box) error {

	// INSERT SQL

	// GET id and store in BoxCache

	return nil
}

package controllers

import (
	"errors"
	"marb.ec/corvi-backend/models"
	"math/rand"
)

var mockBoxes = []*models.Box{
	&models.Box{
		ID:               1,
		Name:             "SQL Statements",
		Category:         mockCategories[0],
		QuestionsToLearn: 2,
		QuestionsTotal:   2,
		QuestionsLearned: 0,
	},
	&models.Box{
		ID:               2,
		Name:             "English - Kitchen",
		Category:         mockCategories[1],
		QuestionsToLearn: 1,
		QuestionsTotal:   1,
		QuestionsLearned: 0,
	},
	&models.Box{
		ID:               3,
		Name:             "French - Cuisine",
		Category:         mockCategories[1],
		QuestionsToLearn: 0,
		QuestionsTotal:   1,
		QuestionsLearned: 1,
	},
}

var BoxControllerSingleton *BoxController

type BoxController struct {
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
	mockBoxes[0].Name = mockBoxes[0].Name + "A"

	// Save in BaxCache

	return mockBoxes, nil
}

func (c *BoxController) LoadBox(id uint) (*models.Box, error) {

	for _, box := range mockBoxes {
		if box.ID == id {
			box.Name = box.Name + "D"
			return box, nil
		}
	}

	// Load Box SQL

	// Store in Cache
	//BoxCache[box.ID] = box

	return nil, errors.New("Box not found.")
}

func (c *BoxController) refreshBox(box *models.Box) error {

	// Subqueries for numbers

	return nil
}

func (c *BoxController) LoadQuestions(id uint) ([]*models.Question, error) {

	_, err := c.LoadBox(id)
	if err != nil {
		return nil, err
	}

	questions := []*models.Question{}
	for _, q := range mockQuestions {
		if q.Box.ID == id {
			questions = append(questions, q)
		}
	}

	return questions, nil
}

func (c *BoxController) GetQuestionToLearn(id uint) (*models.Question, error) {
	_, err := c.LoadBox(id)
	if err != nil {
		return nil, err
	}

	questions, err := c.LoadQuestions(id)
	if err != nil {
		return nil, err
	}

	index := rand.Intn(len(questions) + 1)
	if index < len(questions) {
		return questions[index], nil
	}

	/*

		// Load more questions if heap is nil.
		if box.QuestionsToLearn.Length() == 0 {
			c.loadQuestionsToLearn(box)
		}

		// Get next question in heap. Might be nil if no questions are remaining for this day
		question := box.QuestionsToLearn.Min()

		if question == nil {
			return nil, errors.New("No question to learn for this box.")
		}

		return question, nil */
	return nil, nil
}

func (c *BoxController) loadQuestionsToLearn(b *models.Box) {

	// SELECT * FROM Questions WHERE BoxID = b.ID ORDER BY Next DESC LIMIT 20

	// Skip if question.Next.After(time.Now())
	// TODO(mjb): Update to give questions for whole day
	// TODO(mjb): Add possibility to let user  learn ahead

	//b.questionsToLearn.Add()

}

func (c *BoxController) UpdateBox(b *models.Box) error {
	return nil
}

func (c *BoxController) Insert(b *models.Box) error {

	// INSERT SQL

	// GET id and store in BoxCache

	return nil
}

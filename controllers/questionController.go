package controllers

import (
	"errors"
	"fmt"
	"marb.ec/corvi-backend/models"
	"marb.ec/maf/events"
	"time"
)

var mockQuestions = []*models.Question{
	&models.Question{
		ID:                1,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Box:               mockBoxes[0],
		Next:              time.Now(),
		CorrectlyAnswered: 0,
	},
	&models.Question{
		ID:                2,
		Question:          "Insert Statement?",
		Answer:            "INSERT INTO table (key, key) VALUES (value, value)",
		BoxID:             1,
		Box:               mockBoxes[0],
		Next:              time.Now(),
		CorrectlyAnswered: 1,
	},
	&models.Question{
		ID:                3,
		Question:          "Küche",
		Answer:            "Kitchen",
		BoxID:             2,
		Box:               mockBoxes[1],
		Next:              time.Now(),
		CorrectlyAnswered: 0,
	},
	&models.Question{
		ID:                4,
		Question:          "Küche",
		Answer:            "Cuisine",
		BoxID:             3,
		Box:               mockBoxes[2],
		Next:              time.Now(),
		CorrectlyAnswered: 3,
	},
}

var mockAnswers = []*models.LearnUnit{}

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
	return mockQuestions, nil
}

func (c *QuestionController) LoadQuestion(id uint) (*models.Question, error) {

	for _, q := range mockQuestions {
		if q.ID == id {
			return q, nil
		}
	}

	return nil, errors.New("Question not found.")
}

func (c *QuestionController) GiveCorrectAnswer(id uint) error {
	question, err := c.LoadQuestion(id)
	if err != nil {
		return err
	}

	question.CorrectlyAnswered++
	question.CalculateNext()

	c.saveLearnUnit(question, true)

	// Save Question
	err = c.UpdateQuestion(question)
	if err != nil {
		return err
	}

	// TODO(mjb): Rethink architecture here
	box, err := BoxControllerInstance().LoadBox(question.BoxID)
	if err != nil {
		return err
	}

	BoxControllerInstance().removeQuestionFromHeap(box, question)
	BoxControllerInstance().refreshBox(box)

	return nil
}

func (c *QuestionController) GiveWrongAnswer(id uint) error {

	question, err := c.LoadQuestion(id)
	if err != nil {
		return err
	}
	question.CorrectlyAnswered = 0
	question.CalculateNext()

	c.saveLearnUnit(question, false)

	// Save Question
	err = c.UpdateQuestion(question)
	if err != nil {
		return err
	}

	// TODO(mjb): Rethink architecture here
	box, err := BoxControllerInstance().LoadBox(question.BoxID)
	if err != nil {
		return err
	}

	if relearnUntilAccomplished {
		BoxControllerInstance().reAddQuestionFromHeap(box, question)
	} else {
		BoxControllerInstance().removeQuestionFromHeap(box, question)
	}
	BoxControllerInstance().refreshBox(box)

	return nil
}

// TODO(mjb): Move to LearnUnit Controller
func (c *QuestionController) saveLearnUnit(q *models.Question, correct bool) {
	unit := models.NewLearnUnit(q, correct)
	mockAnswers = append(mockAnswers, unit)
}

func (c *QuestionController) UpdateQuestion(q *models.Question) error {
	events.Events().Publish(events.Topic(fmt.Sprintf("question-%d", q.ID)), c)
	return nil
}

func (c *QuestionController) InsertQuestion(q *models.Question) error {
	BoxControllerInstance().refreshBox(q.Box)
	events.Events().Publish(events.Topic("questions"), c)

	// Update Box
	return nil
}

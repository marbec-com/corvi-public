package controllers

import (
	"errors"
	"fmt"
	"marb.ec/corvi-backend/models"
	"marb.ec/maf/events"
	"time"
)

var mockBoxes = []*models.Box{
	&models.Box{
		ID:               1,
		Name:             "SQL Statements",
		CategoryID:       1,
		Category:         mockCategories[0],
		QuestionsToLearn: 2,
		QuestionsTotal:   2,
		QuestionsLearned: 0,
		QuestionHeap:     models.NewQuestionHeap(),
	},
	&models.Box{
		ID:               2,
		Name:             "English - Kitchen",
		CategoryID:       2,
		Category:         mockCategories[1],
		QuestionsToLearn: 1,
		QuestionsTotal:   1,
		QuestionsLearned: 0,
		QuestionHeap:     models.NewQuestionHeap(),
	},
	&models.Box{
		ID:               3,
		Name:             "French - Cuisine",
		CategoryID:       2,
		Category:         mockCategories[1],
		QuestionsToLearn: 0,
		QuestionsTotal:   1,
		QuestionsLearned: 1,
		QuestionHeap:     models.NewQuestionHeap(),
	},
}

// TODO(mjb): Introduce HEAP cache

// TODO(mjb): Replace with dynamic settings variable
const (
	maxToLearn               uint = 10
	relearnUntilAccomplished bool = true
)

var BoxControllerSingleton *BoxController

type BoxController struct {
}

func BoxControllerInstance() *BoxController {
	if BoxControllerSingleton == nil {
		BoxControllerSingleton = NewBoxController()
		for _, box := range mockBoxes {
			BoxControllerSingleton.loadQuestionsToLearn(box)
			BoxControllerSingleton.refreshBox(box)
		}
	}
	return BoxControllerSingleton
}

func NewBoxController() *BoxController {
	return &BoxController{}
}

func (c *BoxController) LoadBoxes() ([]*models.Box, error) {
	return mockBoxes, nil
}

func (c *BoxController) LoadBox(id uint) (*models.Box, error) {

	for _, box := range mockBoxes {
		if box.ID == id {
			return box, nil
		}
	}

	return nil, errors.New("Box not found.")
}

func (c *BoxController) refreshBox(box *models.Box) error {

	// TODO(mjb): Replace with Database query

	box.QuestionsToLearn = uint(box.QuestionHeap.Length())
	var learned, total uint

	for _, question := range mockQuestions {
		if question.BoxID != box.ID {
			continue
		}
		total++
		if question.CorrectlyAnswered > 0 {
			learned++
		}
	}

	box.QuestionsLearned = learned
	box.QuestionsTotal = total

	events.Events().Publish(events.Topic(fmt.Sprintf("box-%d", box.ID)), c)

	return nil
}

func (c *BoxController) LoadQuestions(id uint) ([]*models.Question, error) {

	_, err := c.LoadBox(id)
	if err != nil {
		return nil, err
	}

	questions := []*models.Question{}
	for _, q := range mockQuestions {
		if q.BoxID == id {
			questions = append(questions, q)
		}
	}

	return questions, nil
}

func (c *BoxController) removeQuestionFromHeap(box *models.Box, question *models.Question) {
	box.QuestionHeap.Lock()
	if box.QuestionHeap.Peek() == question {
		box.QuestionHeap.Min()
	}
	box.QuestionHeap.Unlock()
}

func (c *BoxController) reAddQuestionFromHeap(box *models.Box, question *models.Question) {
	box.QuestionHeap.Lock()
	if box.QuestionHeap.Peek() == question {
		box.QuestionHeap.Add(box.QuestionHeap.Min())
	}
	box.QuestionHeap.Unlock()
}

func (c *BoxController) GetQuestionToLearn(id uint) (*models.Question, error) {
	box, err := c.LoadBox(id)
	if err != nil {
		return nil, err
	}

	if box.QuestionHeap.Length() == 0 {
		c.loadQuestionsToLearn(box)
	}

	fmt.Println(box.QuestionHeap)

	return box.QuestionHeap.Peek(), nil
}

func (c *BoxController) getBeginningOfNextDay() time.Time {
	today := time.Now()
	today = today.Add(24 * time.Hour)
	year, month, day := today.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, today.Location())
}

func (c *BoxController) loadQuestionsToLearn(b *models.Box) {
	capacity := c.getCapacity(b)
	if capacity < 0 {
		return
	}

	/*
		 SELECT *
		 FROM questions
		 WHERE boxID = b.ID
		 AND Next >= BOD
		 AND questionID NOT IN (
			 SELECT questionID
			 FROM learnunit
			 WHERE boxID = b.ID
			 AND time = today
		 )

		 -> capacity

	*/

	tomorrow := c.getBeginningOfNextDay()

	set := make(map[*models.Question]bool)

	// Add all questions of that box
	for _, question := range mockQuestions {
		if question.BoxID != b.ID {
			continue
		}
		set[question] = true
	}

	// Mark all questions that were arleady answered today
	yt, mt, dt := time.Now().Date()
	for _, unit := range mockAnswers {
		if unit.BoxID != b.ID {
			continue
		}
		y, m, d := unit.Time.Date()
		if y == yt && m == mt && d == dt {
			set[unit.Question] = false
		}
	}

	// Only add unmarked question that are due
	for question, a := range set {
		if capacity < 0 {
			return
		}

		if question.Next.Before(tomorrow) && a {
			b.QuestionHeap.Add(question)
			capacity--
			fmt.Println(capacity)
		}
	}

}

func (c *BoxController) getCapacity(b *models.Box) uint {
	// TODO(mjb): Update to SQL query to count correct objects
	capacity := maxToLearn
	yt, mt, dt := time.Now().Date()
	for _, unit := range mockAnswers {
		if unit.BoxID != b.ID {
			continue
		}
		y, m, d := unit.Time.Date()
		if y == yt && m == mt && d == dt {
			capacity--
		}
	}

	if capacity < 0 {
		return 0
	}
	return uint(capacity)
}

func (c *BoxController) UpdateBox(b *models.Box) error {
	events.Events().Publish(events.Topic(fmt.Sprintf("box-%d", b.ID)), c)
	return nil
}

func (c *BoxController) AddBox(b *models.Box) error {
	events.Events().Publish(events.Topic("boxes"), c)
	return nil
}

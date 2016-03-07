package controllers

import (
	"errors"
	"marb.ec/corvi-backend/models"
	"marb.ec/maf/events"
	"time"
)

var mockQuestions = []*models.Question{
	&models.Question{
		ID:                1,
		Question:          "A Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                3,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 5,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                4,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 20,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                5,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                6,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                7,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                8,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                9,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                10,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                11,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                12,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                13,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                14,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                15,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                16,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                17,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                18,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                19,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                20,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                21,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                22,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                23,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                24,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                25,
		Question:          "Update Statement?",
		Answer:            "UPDATE table SET key = value",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                2,
		Question:          "Insert Statement?",
		Answer:            "INSERT INTO table (key, key) VALUES (value, value)",
		BoxID:             1,
		Next:              time.Now(),
		CorrectlyAnswered: 1,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                26,
		Question:          "Küche",
		Answer:            "Kitchen",
		BoxID:             2,
		Next:              time.Now(),
		CorrectlyAnswered: 0,
		CreatedAt:         time.Now(),
	},
	&models.Question{
		ID:                27,
		Question:          "Küche",
		Answer:            "Cuisine",
		BoxID:             3,
		Next:              time.Now(),
		CorrectlyAnswered: 3,
		CreatedAt:         time.Now(),
	},
}

var mockQuestionsID uint = 28

var mockAnswers = []*models.LearnUnit{}

var QuestionControllerSingleton *QuestionController

func QuestionCtrl() *QuestionController {
	return QuestionControllerSingleton
}

type QuestionController struct {
	db *DBController
}

func NewQuestionController(db *DBController) *QuestionController {
	c := &QuestionController{
		db: db,
	}
	return c
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

	if question.CorrectlyAnswered == 0 {
		// Answer correct for a previously unlearned question
		c.saveLearnUnit(question, true, false)
	} else {
		// Answer correct for a previously learned question
		c.saveLearnUnit(question, true, true)
	}

	question.CorrectlyAnswered++
	question.CalculateNext()

	// Save Question
	err = c.UpdateQuestion(question.ID, question)
	if err != nil {
		return err
	}

	// TODO(mjb): Rethink architecture here
	box, err := BoxCtrl().LoadBox(question.BoxID)
	if err != nil {
		return err
	}

	BoxCtrl().removeQuestionFromHeap(box, question)
	BoxCtrl().refreshBox(box)
	events.Events().Publish(events.Topic("stats"), c)

	return nil
}

func (c *QuestionController) GiveWrongAnswer(id uint) error {

	question, err := c.LoadQuestion(id)
	if err != nil {
		return err
	}

	if question.CorrectlyAnswered == 0 {
		// Answer false for a previously unlearned question
		c.saveLearnUnit(question, false, false)
	} else {
		// Answer false for a previously learned question
		c.saveLearnUnit(question, false, true)
	}

	question.CorrectlyAnswered = 0
	question.CalculateNext()

	// Save Question
	err = c.UpdateQuestion(question.ID, question)
	if err != nil {
		return err
	}

	// TODO(mjb): Rethink architecture here
	box, err := BoxCtrl().LoadBox(question.BoxID)
	if err != nil {
		return err
	}

	if SettingsCtrl().Get().RelearnUntilAccomplished {
		BoxCtrl().reAddQuestionFromHeap(box, question)
	} else {
		BoxCtrl().removeQuestionFromHeap(box, question)
	}
	BoxCtrl().refreshBox(box)
	events.Events().Publish(events.Topic("stats"), c)

	return nil
}

// TODO(mjb): Move to LearnUnit Controller
func (c *QuestionController) saveLearnUnit(q *models.Question, correct, prev bool) {
	unit := models.NewLearnUnit(q.ID, q.BoxID, correct, prev)
	mockAnswers = append(mockAnswers, unit)
}

func (c *QuestionController) UpdateQuestion(qID uint, question *models.Question) error {

	for k, q := range mockQuestions {
		if q.ID == qID {
			prevBox := mockQuestions[k].BoxID
			mockQuestions[k] = question

			// Question might have been moved
			if prevBox != question.BoxID {
				// Refresh both boxes
				if prevBoxInstance, err := BoxCtrl().LoadBox(question.BoxID); err == nil {
					BoxCtrl().rebuildQuestionHeap(prevBoxInstance)
					BoxCtrl().refreshBox(prevBoxInstance)
				}
				if curBoxInstance, err := BoxCtrl().LoadBox(prevBox); err == nil {
					BoxCtrl().rebuildQuestionHeap(curBoxInstance)
					BoxCtrl().refreshBox(curBoxInstance)
				}
				events.Events().Publish(events.Topic("boxes"), c)
			}
			events.Events().Publish(events.Topic("questions"), c)

			return nil
		}
	}

	return errors.New("Question to update was not found.")
}

func (c *QuestionController) AddQuestion(q *models.Question) (*models.Question, error) {
	box, err := BoxCtrl().LoadBox(q.BoxID)
	if err != nil {
		return nil, errors.New("Box for this question does not exist.")
	}

	// Insert question
	q.ID = mockQuestionsID
	mockQuestionsID++
	mockQuestions = append(mockQuestions, q)

	// Rebuild heap and refresh stats
	BoxCtrl().rebuildQuestionHeap(box)
	BoxCtrl().refreshBox(box)

	events.Events().Publish(events.Topic("questions"), c)
	events.Events().Publish(events.Topic("stats"), c)

	return q, nil
}

func (c *QuestionController) DeleteQuestion(qID uint) error {

	// TODO(mjb): Remove answers from that question

	for k, q := range mockQuestions {
		if q.ID == qID {
			boxID := q.BoxID
			mockQuestions, mockQuestions[len(mockQuestions)-1] = append(mockQuestions[:k], mockQuestions[k+1:]...), nil

			// Refresh box of deleted question
			box, _ := BoxCtrl().LoadBox(boxID)
			BoxCtrl().rebuildQuestionHeap(box)
			BoxCtrl().refreshBox(box)

			events.Events().Publish(events.Topic("questions"), c)
			events.Events().Publish(events.Topic("stats"), c)

			return nil
		}
	}

	return errors.New("Question not found.")
}

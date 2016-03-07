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
		QuestionsToLearn: 2,
		QuestionsTotal:   2,
		QuestionsLearned: 0,
		QuestionHeap:     models.NewQuestionHeap(),
		CreatedAt:        time.Now(),
	},
	&models.Box{
		ID:               2,
		Name:             "English - Kitchen",
		CategoryID:       2,
		QuestionsToLearn: 1,
		QuestionsTotal:   1,
		QuestionsLearned: 0,
		QuestionHeap:     models.NewQuestionHeap(),
		CreatedAt:        time.Now(),
	},
	&models.Box{
		ID:               3,
		Name:             "French - Cuisine",
		CategoryID:       2,
		QuestionsToLearn: 0,
		QuestionsTotal:   1,
		QuestionsLearned: 1,
		QuestionHeap:     models.NewQuestionHeap(),
		CreatedAt:        time.Now(),
	},
}

var mockBoxesID uint = 4

// TODO(mjb): Introduce HEAP cache

var BoxControllerSingleton *BoxController

func BoxCtrl() *BoxController {
	return BoxControllerSingleton
}

type BoxController struct {
	db *DBController
}

func NewBoxController(db *DBController) (*BoxController, error) {
	b := &BoxController{
		db: db,
	}
	err := b.createTables()
	if err != nil {
		return nil, err
	}
	/* for _, box := range mockBoxes {
		b.loadQuestionsToLearn(box)
		b.refreshBox(box)
	} */
	return b, nil
}

func (c *BoxController) createTables() error {

	// Create table, only if it not already exists
	// Includes foreign key constraint to Category table. We are not allowed to delete a Category that still has Boxes assigned.
	sql := "CREATE TABLE IF NOT EXISTS Box (ID INTEGER PRIMARY KEY ASC NOT NULL, Name VARCHAR (255) NOT NULL, Description TEXT NOT NULL, CategoryID INTEGER REFERENCES Category (ID) NOT NULL, CreatedAt DATETIME NOT NULL);"
	_, err := c.db.Connection().Exec(sql)

	return err

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

func (c *BoxController) rebuildQuestionHeap(box *models.Box) {
	box.QuestionHeap.Clear()
	c.loadQuestionsToLearn(box)
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
	if capacity <= 0 {
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

	set := make(map[uint]*models.Question)

	// Add all questions of that box
	for _, question := range mockQuestions {
		if question.BoxID != b.ID {
			continue
		}
		set[question.ID] = question
	}

	// Set all questions that were arleady answered today to nil
	yt, mt, dt := time.Now().Date()
	for _, unit := range mockAnswers {
		if unit.BoxID != b.ID {
			continue
		}
		y, m, d := unit.CreatedAt.Date()
		if y == yt && m == mt && d == dt {
			set[unit.QuestionID] = nil
		}
	}

	// Only add not nil question that are due
	for _, question := range set {
		if capacity <= 0 {
			return
		}
		if question != nil && question.Next.Before(tomorrow) {
			b.QuestionHeap.Add(question)
			capacity--
		}
	}

}

func (c *BoxController) getCapacity(b *models.Box) uint {
	// TODO(mjb): Update to SQL query to count correct objects
	capacity := SettingsCtrl().Get().MaxDailyQuestionsPerBox
	yt, mt, dt := time.Now().Date()
	for _, unit := range mockAnswers {
		if capacity <= 0 {
			return 0
		}
		if unit.BoxID != b.ID {
			continue
		}
		y, m, d := unit.CreatedAt.Date()
		if y == yt && m == mt && d == dt {
			capacity--
		}
	}

	return capacity
}

func (c *BoxController) UpdateBox(boxID uint, box *models.Box) error {

	for k, b := range mockBoxes {
		if b.ID == boxID {
			mockBoxes[k] = box
			// Update Box, Previous cat, new cat
			events.Events().Publish(events.Topic("boxes"), c)
			return nil
		}
	}

	return errors.New("Box to update was not found.")
}

func (c *BoxController) AddBox(box *models.Box) (*models.Box, error) {
	box.ID = mockBoxesID
	mockBoxesID++

	mockBoxes = append(mockBoxes, box)

	events.Events().Publish(events.Topic("boxes"), c)
	events.Events().Publish(events.Topic("stats"), c)

	return box, nil
}

func (c *BoxController) DeleteBox(boxID uint) error {

	// TODO(mjb): Remove answers from all questions

	// Delete all questions of that box, start with highest index so that following indexes do not move
	qIndexes := []int{}
	for k := len(mockQuestions) - 1; k >= 0; k-- {
		if mockQuestions[k].BoxID == boxID {
			qIndexes = append(qIndexes, k)
		}
	}

	for _, i := range qIndexes {
		mockQuestions, mockQuestions[len(mockQuestions)-1] = append(mockQuestions[:i], mockQuestions[i+1:]...), nil
	}

	// Delete box
	for k, b := range mockBoxes {
		if b.ID == boxID {
			mockBoxes, mockBoxes[len(mockBoxes)-1] = append(mockBoxes[:k], mockBoxes[k+1:]...), nil

			events.Events().Publish(events.Topic("questions"), c)
			events.Events().Publish(events.Topic("boxes"), c)
			events.Events().Publish(events.Topic("stats"), c)
			return nil
		}
	}

	return errors.New("Box not found")
}

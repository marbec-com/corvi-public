package controllers

import (
	"log"
	"marb.ec/corvi-backend/models"
	"marb.ec/maf/events"
	"marb.ec/maf/interfaces"
)

type MockSubscriber struct {
	Notifications map[string]uint
}

func NewMockSubscriber(topics []string) *MockSubscriber {
	s := &MockSubscriber{
		Notifications: make(map[string]uint, len(topics)),
	}

	for _, topic := range topics {
		events.Events().Subscribe(events.Topic(topic), s)
	}
	return s
}

func (s *MockSubscriber) Notify(t interfaces.Topic, p interfaces.Publisher) {
	s.Notifications[t.Topic()]++
}

func (s *MockSubscriber) Assert(topic string, count uint) bool {
	c, ok := s.Notifications[topic]
	if !ok {
		return false
	}
	return c == count
}

type MockSettingsService struct {
	settings *models.Settings
}

func NewMockSettingsService(settings *models.Settings) *MockSettingsService {
	return &MockSettingsService{
		settings: settings,
	}
}

func (s *MockSettingsService) Update() error {
	return nil
}

func (s *MockSettingsService) Get() *models.Settings {
	return s.settings
}

type MockBoxController struct {
	BuildHeapCount              int
	RemoveQuestionFromHeapCount int
	ReAddQuestionFromHeapCount  int
}

func (c *MockBoxController) CreateTables() error {
	return nil
}

func (c *MockBoxController) LoadBoxes() ([]*models.Box, error) {
	return []*models.Box{}, nil
}

func (c *MockBoxController) LoadBoxesOfCategory(id uint) ([]*models.Box, error) {
	return []*models.Box{}, nil
}

func (c *MockBoxController) LoadBox(id uint) (*models.Box, error) {
	return nil, nil
}

func (c *MockBoxController) RemoveQuestionFromHeap(id, qID uint) error {
	c.RemoveQuestionFromHeapCount++
	return nil
}

func (c *MockBoxController) ReAddQuestionFromHeap(id, qID uint) error {
	c.ReAddQuestionFromHeapCount++
	return nil
}

func (c *MockBoxController) GetQuestionToLearn(id uint) (*models.Question, error) {
	return nil, nil
}

func (c *MockBoxController) BuildHeap(id uint) error {
	c.BuildHeapCount++
	return nil
}

func (c *MockBoxController) BuildHeaps() error {
	return nil
}

func (c *MockBoxController) UpdateBox(boxID uint, box *models.Box) error {
	return nil
}

func (c *MockBoxController) AddBox(box *models.Box) (*models.Box, error) {
	return nil, nil
}

func (c *MockBoxController) DeleteBox(boxID uint) error {
	return nil
}

func setupTestDBController(path string) DatabaseService {
	c, err := NewSQLiteDBService(path)
	if err != nil {
		log.Fatal("Error in Setup", err)
		return nil
	}
	return c
}

func tearDownTestDBController(db DatabaseService) {
	// Close database connection
	err := db.Close()
	if err != nil {
		log.Fatal("Error in Teardown", err)
		return
	}
	// Remove file
	err = db.Destroy()
	if err != nil {
		log.Fatal("Error in Teardown", err)
		return
	}
}

func insertRawCategories(categories []*models.Category, db DatabaseService) {

	for _, c := range categories {
		sqlStmt := "INSERT INTO Category (Name, CreatedAt) VALUES (?, ?);"
		res, err := db.Connection().Exec(sqlStmt, c.Name, c.CreatedAt)
		if err != nil {
			log.Fatal("Could not insert test data", err)
		}

		// Update objects ID
		newID, err := res.LastInsertId()
		if err != nil {
			log.Fatal("Could not update IDs of test data", err)
		}
		c.ID = uint(newID)
	}

}

func insertRawBoxes(boxes []*models.Box, db DatabaseService) {

	for _, b := range boxes {
		sqlStmt := "INSERT INTO Box (Name, Description, CategoryID, CreatedAt) VALUES (?, ?, ?, ?);"
		res, err := db.Connection().Exec(sqlStmt, b.Name, b.Description, b.CategoryID, b.CreatedAt)
		if err != nil {
			log.Fatal("Could not insert test data", err)
		}

		// Update objects ID
		newID, err := res.LastInsertId()
		if err != nil {
			log.Fatal("Could not update IDs of test data", err)
		}
		b.ID = uint(newID)
	}

}

func insertRawQuestions(questions []*models.Question, db DatabaseService) {

	for _, q := range questions {
		sqlStmt := "INSERT INTO Question (Question, Answer, BoxID, Next, CorrectlyAnswered, CreatedAt) VALUES (?, ?, ?, ?, ?, ?);"
		res, err := db.Connection().Exec(sqlStmt, q.Question, q.Answer, q.BoxID, q.Next, q.CorrectlyAnswered, q.CreatedAt)
		if err != nil {
			log.Fatal("Could not insert test data", err)
		}

		// Update objects ID
		newID, err := res.LastInsertId()
		if err != nil {
			log.Fatal("Could not update IDs of test data", err)
		}
		q.ID = uint(newID)
	}

}

func insertRawLearnUnits(learnUnits []*models.LearnUnit, db DatabaseService) {

	for _, l := range learnUnits {
		sqlStmt := "INSERT INTO LearnUnit(QuestionID, BoxID, Correct, PrevCorrect, CreatedAt) VALUES (?, ?, ?, ?, ?);"
		_, err := db.Connection().Exec(sqlStmt, l.QuestionID, l.BoxID, l.Correct, l.PrevCorrect, l.CreatedAt)
		if err != nil {
			log.Fatal("Could not insert test data", err)
		}
	}

}

func getLearnUnit(qID uint, db DatabaseService) *models.LearnUnit {

	sql := "SELECT QuestionID, BoxID, Correct, PrevCorrect, CreatedAt FROM LearnUnit WHERE QuestionID = ? LIMIT 1;"
	row := db.Connection().QueryRow(sql, qID)

	// Create new Category object
	lu := models.NewLearnUnit()

	// Populate
	err := row.Scan(&lu.QuestionID, &lu.BoxID, &lu.Correct, &lu.PrevCorrect, &lu.CreatedAt)
	if err != nil {
		log.Fatal("Could not retrieve test data", err)
	}

	return lu

}

package controllers

import (
	"fmt"
	"log"
	"marb.ec/corvi-backend/models"
	"testing"
	"time"
)

func setupTestQuestionController(path string, settings *models.Settings) (DatabaseService, *QuestionController, *BoxControllerImpl, *CategoryController) {

	db := setupTestDBController(path)
	s := NewMockSettingsService(settings)

	// Make sure category tables are created
	c, err := NewCategoryController(db, s)
	if err != nil {
		log.Fatal("Error in Setup", err)
		return nil, nil, nil, nil
	}

	// Make sure question tables are created
	q, err := NewQuestionController(db, s)
	if err != nil {
		log.Fatal("Error in Setup", err)
		return nil, nil, nil, nil
	}

	b, err := NewBoxController(db, s)
	if err != nil {
		log.Fatal("Error in Setup", err)
		return nil, nil, nil, nil
	}
	return db, q, b, c

}

func insertQuestionTestData(db DatabaseService) (*models.Question, *models.Question) {

	boxA, boxB := insertBoxTestData(db)

	questionA := models.NewQuestion()
	questionA.Question = "Question A"
	questionA.Answer = "Answer A"
	questionA.BoxID = boxA.ID
	questionA.CreatedAt = time.Now()
	questionA.CalculateNext()
	questionA.CorrectlyAnswered = 10
	questionB := models.NewQuestion()
	questionB.Question = "question B"
	questionB.Answer = "Answer B"
	questionB.BoxID = boxB.ID
	questionB.CreatedAt = time.Now()
	questionB.CalculateNext()

	insertRawQuestions([]*models.Question{questionA, questionB}, db)

	return questionA, questionB

}

func TestQuestionCtrlCreateTables(t *testing.T) {

	// Setup & Teardown
	db := setupTestDBController("test_questionController.db")
	defer tearDownTestDBController(db)

	// Create controller
	questionController := &QuestionController{
		db: db,
	}

	// Execute createTables()
	err := questionController.createTables()
	if err != nil {
		t.Log("Error while executing createTables", err)
		t.Fail()
	}

	// Check SQL
	sqlStmt := "SELECT COUNT(*) FROM sqlite_master WHERE (type = 'table' AND name = 'Question') OR (type = 'table' AND name = 'LearnUnit') OR (type = 'view' AND name = 'QuestionsLearnedToday') OR (type = 'view' AND name = 'QuestionsDue');"
	row := questionController.db.Connection().QueryRow(sqlStmt)
	var count int
	err = row.Scan(&count)
	if err != nil || count != 4 {
		t.Log("Table and view were not created", err, count)
		t.Fail()
	}

}

func TestQuestionCtrlLoadQuestions(t *testing.T) {

	// Setup & Teardown
	db, controller, _, _ := setupTestQuestionController("test_questionController.db", nil)
	defer tearDownTestDBController(db)
	qA, qB := insertQuestionTestData(db)

	// Load Questions
	questions, err := controller.LoadQuestions()
	if err != nil {
		t.Log("LoadQuestions returned an error", err)
		t.Fail()
	}

	// Check length
	if len(questions) != 2 {
		t.Log("Returned array does not have length of 2:", len(questions))
		t.Fail()
	}

	// Compare
	if !((questions[0].Equal(qA) && questions[1].Equal(qB)) || (questions[0].Equal(qB) && questions[1].Equal(qA))) {
		t.Log("Returned questions do not match", questions, qA, qB)
		t.Fail()
	}

}

func TestQuestionCtrlLoadQuestionsOfBox(t *testing.T) {

	// Setup & Teardown
	db, controller, _, _ := setupTestQuestionController("test_questionController.db", nil)
	defer tearDownTestDBController(db)
	qA, _ := insertQuestionTestData(db)

	// Load Questions
	questions, err := controller.LoadQuestionsOfBox(qA.BoxID)
	if err != nil {
		t.Log("LoadQuestions returned an error", err)
		t.Fail()
	}

	// Compare
	if len(questions) != 1 || !questions[0].Equal(qA) {
		t.Log("Returned array does not contain qA", questions)
		t.Fail()
	}

}

func TestQuestionCtrlLoadQuestion(t *testing.T) {

	// Setup & Teardown
	db, controller, _, _ := setupTestQuestionController("test_questionController.db", nil)
	defer tearDownTestDBController(db)
	_, qB := insertQuestionTestData(db)

	// Get second
	question, err := controller.LoadQuestion(qB.ID)
	if err != nil {
		t.Log("LoadQuestion returned an error", err)
		t.Fail()
	}

	// Compare
	if question == nil || !question.Equal(qB) {
		t.Log("Second inserted question does not equal result", qB, question)
		t.Fail()
	}

}

func TestQuestionCtrlUpdateQuestion(t *testing.T) {

	// Setup & Teardown
	db, controller, _, _ := setupTestQuestionController("test_questionController.db", nil)
	mockBoxCtrl := &MockBoxController{}
	controller.boxCtrl = mockBoxCtrl
	defer tearDownTestDBController(db)
	qA, qB := insertQuestionTestData(db)
	sub := NewMockSubscriber([]string{"boxes", "questions"})

	// Manipulate questionB
	qB.Question = "New Question for B"
	qB.Answer = "Revised answer"
	qB.CorrectlyAnswered = 5
	qB.BoxID = qA.BoxID // Change box
	qB.CalculateNext()
	qB.CreatedAt = time.Now()

	// Update
	err := controller.UpdateQuestion(qB.ID, qB)
	if err != nil {
		t.Log("UpdateQuestion returned an error", err)
		t.Fail()
	}

	// Load
	question, err := controller.LoadQuestion(qB.ID)
	if err != nil {
		t.Log("LoadQuestion returned an error", err)
		t.Fail()
	}

	// Compare
	if question == nil || !question.Equal(qB) {
		t.Log("Retrieved question does not equal update", question, qB)
		t.Fail()
	}

	// Check function calls
	if mockBoxCtrl.BuildHeapCount != 2 { // For original and new Box
		t.Log("BuildHeap was not called.", mockBoxCtrl.BuildHeapCount)
		t.Fail()
	}

	// Notifications
	if !sub.Assert("boxes", 1) {
		t.Log("Update notification boxes was not sent.")
		t.Fail()
	}
	if !sub.Assert("questions", 1) {
		t.Log("Update notification questions was not sent.")
		t.Fail()
	}

}

func TestQuestionCtrlAddQuestion(t *testing.T) {

	// Setup & Teardown
	db, controller, _, _ := setupTestQuestionController("test_questionController.db", nil)
	mockBoxCtrl := &MockBoxController{}
	controller.boxCtrl = mockBoxCtrl
	defer tearDownTestDBController(db)
	qA, _ := insertQuestionTestData(db)
	sub := NewMockSubscriber([]string{"stats", "questions", fmt.Sprintf("box-%d", qA.BoxID)})

	// Create QuestionC
	qC := models.NewQuestion()
	qC.ID = 0
	qC.Question = "New Question C"
	qC.Answer = "Answer of Question C"
	qC.CorrectlyAnswered = 5
	qC.BoxID = qA.BoxID // Change box
	qC.CalculateNext()
	qC.CreatedAt = time.Now()

	// Insert
	qC, err := controller.AddQuestion(qC)
	if err != nil {
		t.Log("AddQuestion returned an error", err)
		t.Fail()
	}

	// Check ID
	if qC.ID <= 0 {
		t.Log("ID field of inserted question was not updated", qC.ID)
		t.Fail()
	}

	// Load
	retrieveQ, err := controller.LoadQuestion(qC.ID)
	if err != nil {
		t.Log("LoadQuestion returned an error", err)
		t.Fail()
	}

	// Compare
	if retrieveQ == nil || !retrieveQ.Equal(qC) {
		t.Log("Retrieved question does not equal inserted", retrieveQ, qC)
		t.Fail()
	}

	// Check function calls
	if mockBoxCtrl.BuildHeapCount != 1 {
		t.Log("BuildHeap was not called.", mockBoxCtrl.BuildHeapCount)
		t.Fail()
	}

	// Notifications
	if !sub.Assert(fmt.Sprintf("box-%d", qA.BoxID), 1) {
		t.Log("Insert notification boxes was not sent.")
		t.Fail()
	}
	if !sub.Assert("questions", 1) {
		t.Log("Insert notification boxes was not sent.")
		t.Fail()
	}
	if !sub.Assert("stats", 1) {
		t.Log("Insert notification stats was not sent.")
		t.Fail()
	}

}

func TestQuestionCtrlDeleteQuestion(t *testing.T) {

	// Setup & Teardown
	db, controller, _, _ := setupTestQuestionController("test_questionController.db", nil)
	mockBoxCtrl := &MockBoxController{}
	controller.boxCtrl = mockBoxCtrl
	defer tearDownTestDBController(db)
	_, qB := insertQuestionTestData(db)
	sub := NewMockSubscriber([]string{"stats", "questions", fmt.Sprintf("box-%d", qB.BoxID)})

	// Insert LearnUnits for Question
	lA := models.NewLearnUnit()
	lA.QuestionID = qB.ID
	lA.BoxID = qB.BoxID
	lA.Correct = true
	lA.PrevCorrect = false
	lB := models.NewLearnUnit()
	lB.QuestionID = qB.ID
	lB.BoxID = qB.BoxID
	lB.Correct = false
	lB.PrevCorrect = false
	insertRawLearnUnits([]*models.LearnUnit{lA, lB}, db)

	// Delete
	err := controller.DeleteQuestion(qB.ID)
	if err != nil {
		t.Log("DeleteQuestion returned an error", err)
		t.Fail()
	}

	// Check if deleted question is still there
	question, err := controller.LoadQuestion(qB.ID)
	if question != nil || err == nil {
		t.Log("Question was not deleted", question)
		t.Fail()
	}

	// Check function calls
	if mockBoxCtrl.BuildHeapCount != 1 {
		t.Log("BuildHeap was not called.", mockBoxCtrl.BuildHeapCount)
		t.Fail()
	}

	// Check if there are sill LearnUnits for that question
	sqlStmt := "SELECT COUNT(*) FROM LearnUnit WHERE QuestionID = ?;"
	row := controller.db.Connection().QueryRow(sqlStmt, qB.ID)
	var count int
	err = row.Scan(&count)
	if err != nil {
		t.Log("Error while counting LearnUnits", err)
		t.Fail()
	}
	if count != 0 {
		t.Log("LearnUnit of question were not deleted", count)
		t.Fail()
	}

	// Notifications
	if !sub.Assert(fmt.Sprintf("box-%d", qB.BoxID), 1) {
		t.Log("Insert notification boxes was not sent.")
		t.Fail()
	}
	if !sub.Assert("questions", 1) {
		t.Log("Insert notification boxes was not sent.")
		t.Fail()
	}
	if !sub.Assert("stats", 1) {
		t.Log("Insert notification stats was not sent.")
		t.Fail()
	}

}

func TestQuestionCtrlGiveAnswerCorrect(t *testing.T) {

	// Setup & Teardown
	settings := models.NewSettings()
	settings.RelearnUntilAccomplished = false
	db, controller, _, _ := setupTestQuestionController("test_questionController.db", settings)
	mockBoxCtrl := &MockBoxController{}
	controller.boxCtrl = mockBoxCtrl
	defer tearDownTestDBController(db)
	qA, _ := insertQuestionTestData(db)
	sub := NewMockSubscriber([]string{"stats", "questions", fmt.Sprintf("box-%d", qA.BoxID)})

	// Create Expected LearnUnit
	lu := models.NewLearnUnit()
	lu.QuestionID = qA.ID
	lu.BoxID = qA.BoxID
	lu.Correct = true
	lu.PrevCorrect = true // qA.CorrectlyAnswered = 10

	// Update Question to expected -> CalculateNext
	qA.CorrectlyAnswered++
	qA.CalculateNext()

	// Give correct answer
	err := controller.GiveAnswer(qA.ID, true)
	if err != nil {
		t.Log("GiveAnswer returned an error", err)
		t.Fail()
	}

	// Load Question
	retrievedQ, err := controller.LoadQuestion(qA.ID)
	if err != nil {
		t.Log("LoadQuestion returned an error", err)
		t.Fail()
	}

	// Load Learn Unit
	retrievedLu := getLearnUnit(qA.ID, db)

	// Compare
	if retrievedQ == nil || !retrievedQ.Equal(qA) {
		t.Log("Retrieved question does not equal inserted", retrievedQ, qA)
		t.Fail()
	}

	// Ignore Created At
	retrievedLu.CreatedAt = lu.CreatedAt

	if retrievedLu == nil || !retrievedLu.Equal(lu) {
		t.Log("Retrieved learnUnit does not equal inserted", retrievedLu, lu)
		t.Fail()
	}

	// Check function calls
	if mockBoxCtrl.RemoveQuestionFromHeapCount != 1 {
		t.Log("ReAddQuestionFromHeap was not called.", mockBoxCtrl.RemoveQuestionFromHeapCount)
		t.Fail()
	}

	// Notifications
	if !sub.Assert(fmt.Sprintf("box-%d", qA.BoxID), 1) {
		t.Log("Insert notification boxes was not sent.")
		t.Fail()
	}
	if !sub.Assert("questions", 1) {
		t.Log("Insert notification boxes was not sent.")
		t.Fail()
	}
	if !sub.Assert("stats", 1) {
		t.Log("Insert notification stats was not sent.")
		t.Fail()
	}
}

func TestQuestionCtrlGiveAnswerWrong(t *testing.T) {

	// Setup & Teardown
	settings := models.NewSettings()
	settings.RelearnUntilAccomplished = true
	db, controller, _, _ := setupTestQuestionController("test_questionController.db", settings)
	mockBoxCtrl := &MockBoxController{}
	controller.boxCtrl = mockBoxCtrl
	defer tearDownTestDBController(db)
	qA, _ := insertQuestionTestData(db)
	sub := NewMockSubscriber([]string{"stats", "questions", fmt.Sprintf("box-%d", qA.BoxID)})

	// Create Expected LearnUnit
	lu := models.NewLearnUnit()
	lu.QuestionID = qA.ID
	lu.BoxID = qA.BoxID
	lu.Correct = false
	lu.PrevCorrect = true // qA.CorrectlyAnswered = 10

	// Update Question to expected -> CalculateNext
	qA.CorrectlyAnswered = 0
	qA.CalculateNext()

	// Give correct answer
	err := controller.GiveAnswer(qA.ID, false)
	if err != nil {
		t.Log("GiveAnswer returned an error", err)
		t.Fail()
	}

	// Load Question
	retrievedQ, err := controller.LoadQuestion(qA.ID)
	if err != nil {
		t.Log("LoadQuestion returned an error", err)
		t.Fail()
	}

	// Load Learn Unit
	retrievedLu := getLearnUnit(qA.ID, db)

	// Compare
	if retrievedQ == nil || !retrievedQ.Equal(qA) {
		t.Log("Retrieved question does not equal inserted", retrievedQ, qA)
		t.Fail()
	}

	// Ignore Created At
	retrievedLu.CreatedAt = lu.CreatedAt

	if retrievedLu == nil || !retrievedLu.Equal(lu) {
		t.Log("Retrieved learnUnit does not equal inserted", retrievedLu, lu)
		t.Fail()
	}

	// Check function calls
	if mockBoxCtrl.ReAddQuestionFromHeapCount != 1 {
		t.Log("RemoveQuestionFromHeap was not called.", mockBoxCtrl.ReAddQuestionFromHeapCount)
		t.Fail()
	}

	// Notifications
	if !sub.Assert(fmt.Sprintf("box-%d", qA.BoxID), 1) {
		t.Log("Insert notification boxes was not sent.")
		t.Fail()
	}
	if !sub.Assert("questions", 1) {
		t.Log("Insert notification boxes was not sent.")
		t.Fail()
	}
	if !sub.Assert("stats", 1) {
		t.Log("Insert notification stats was not sent.")
		t.Fail()
	}

}

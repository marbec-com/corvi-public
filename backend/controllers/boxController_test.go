package controllers

import (
	"fmt"
	"log"
	"marb.ec/corvi-backend/models"
	"testing"
	"time"
)

func setupTestBoxController(path string) (*BoxController, *CategoryController, *QuestionController) {

	db := setupTestDBController(path)

	// Make sure category tables are created
	c, err := NewCategoryController(db)
	if err != nil {
		log.Fatal("Error in Setup", err)
		return nil, nil, nil
	}

	// Make sure question tables are created
	q, err := NewQuestionController(db)
	if err != nil {
		log.Fatal("Error in Setup", err)
		return nil, nil, nil
	}

	b, err := NewBoxController(db)
	if err != nil {
		log.Fatal("Error in Setup", err)
		return nil, nil, nil
	}
	return b, c, q

}

func insertBoxTestData(boxCtrl *BoxController, catCtrl *CategoryController) (*models.Box, *models.Box) {

	catA, catB := insertCategoryTestData(catCtrl)

	now := time.Now()
	boxA := models.NewBox()
	boxA.ID = 5
	boxA.Name = "Box A"
	boxA.Description = "This is Box A"
	boxA.CategoryID = catA.ID
	boxA.CreatedAt = now
	boxB := models.NewBox()
	boxB.ID = 6
	boxB.Name = "Box B"
	boxB.Description = "This is another Box"
	boxB.CategoryID = catB.ID
	boxB.CreatedAt = now

	// SQL INSERT two boxes
	sqlStmt := "INSERT INTO Box (ID, Name, Description, CategoryID, CreatedAt) VALUES (?,?,?,?,?), (?,?,?,?,?);"
	_, err := boxCtrl.db.Connection().Exec(sqlStmt, boxA.ID, boxA.Name, boxA.Description, boxA.CategoryID, boxA.CreatedAt, boxB.ID, boxB.Name, boxB.Description, boxB.CategoryID, boxB.CreatedAt)
	if err != nil {
		log.Fatal("Could not insert test data", err)
	}

	return boxA, boxB

}

func insertBoxTestDataQuestionsForBox(boxID uint, qCtrl *QuestionController) (*models.Question, *models.Question) {

	questionA := models.NewQuestion()
	questionA.Question = "QuestionA"
	questionA.Answer = "AnswerA"
	questionA.BoxID = boxID
	questionA.CreatedAt = time.Now()
	questionA.CalculateNext()
	questionB := models.NewQuestion()
	questionB.Question = "questionB"
	questionB.Answer = "AnswerB"
	questionB.BoxID = boxID
	questionB.CreatedAt = time.Now()
	questionB.CalculateNext()

	// SQL INSERT two questions
	sqlStmt := "INSERT INTO Question (Question, Answer, BoxID, Next, CorrectlyAnswered, CreatedAt) VALUES (?,?,?,?,?,?), (?,?,?,?,?,?);"
	_, err := qCtrl.db.Connection().Exec(sqlStmt, questionA.Question, questionA.Answer, questionA.BoxID, questionA.Next, questionA.CorrectlyAnswered, questionA.CreatedAt, questionB.Question, questionB.Answer, questionB.BoxID, questionB.Next, questionB.CorrectlyAnswered, questionB.CreatedAt)
	if err != nil {
		log.Fatal("Could not insert test data", err)
	}

	return questionA, questionB

}

func TestBoxCtrlCreateTables(t *testing.T) {

	// Setup & Teardown
	db := setupTestDBController("test_boxController.db")
	defer tearDownTestDBController(db)

	// Create controller
	boxController := &BoxController{
		db: db,
	}

	// Execute createTables()
	err := boxController.createTables()
	if err != nil {
		t.Log("Error while executing createTables", err)
		t.Fail()
	}

	// Check SQL
	sqlStmt := "SELECT COUNT(*) FROM sqlite_master WHERE (type = 'table' AND name = 'Box') OR (type = 'view' AND name = 'BoxWithMeta')"
	row := boxController.db.Connection().QueryRow(sqlStmt)
	var count int
	err = row.Scan(&count)
	if err != nil || count != 2 {
		t.Log("Table and view were not created", err, count)
		t.Fail()
	}

}

func TestBoxCtrlLoadBoxes(t *testing.T) {

	// Setup & Teardown
	controller, catCtrl, _ := setupTestBoxController("test_boxController.db")
	defer tearDownTestDBController(controller.db)
	boxA, boxB := insertBoxTestData(controller, catCtrl)

	// Create Heap for single Box
	heap := models.NewQuestionHeap()
	heap.Add(models.NewQuestion())
	heap.Add(models.NewQuestion())
	controller.heapCache[boxA.ID] = heap

	// Expect QuestionsToLearn to be 2
	boxA.QuestionsToLearn = uint(heap.Length())

	// Load Boxes
	boxes, err := controller.LoadBoxes()
	if err != nil {
		t.Log("LoadBoxes returned an error", err)
		t.Fail()
	}

	// Compare
	if len(boxes) != 2 {
		t.Log("Returned array does not have length of 2:", len(boxes))
		t.Fail()
	}

	if !((boxes[0].Equal(boxA) && boxes[1].Equal(boxB)) || (boxes[0].Equal(boxB) && boxes[1].Equal(boxA))) {
		t.Log("Returned boxes do not match", boxes, boxA, boxB)
		t.Fail()
	}

}

func TestBoxCtrlLoadBoxesOfCategory(t *testing.T) {

	// Setup & Teardown
	controller, catCtrl, _ := setupTestBoxController("test_boxController.db")
	defer tearDownTestDBController(controller.db)
	_, boxB := insertBoxTestData(controller, catCtrl)

	// Load Boxes
	boxes, err := controller.LoadBoxesOfCategory(boxB.CategoryID)
	if err != nil {
		t.Log("LoadBoxes returned an error", err)
		t.Fail()
	}

	if len(boxes) != 1 || !boxes[0].Equal(boxB) {
		t.Log("Returned array does not contain boxB", boxes)
		t.Fail()
	}

}

func TestBoxCtrlLoadBox(t *testing.T) {

	// Setup & Teardown
	controller, catCtrl, _ := setupTestBoxController("test_boxController.db")
	defer tearDownTestDBController(controller.db)
	_, boxB := insertBoxTestData(controller, catCtrl)

	// Get second
	box, err := controller.LoadBox(boxB.ID)
	if err != nil {
		t.Log("LoadBox returned an error", err)
		t.Fail()
	}

	// Compare
	if box == nil || !box.Equal(boxB) {
		t.Log("Second inserted box does not equal result", boxB, box)
		t.Fail()
	}

}

func TestBoxCtrlUpdateBox(t *testing.T) {

	// Setup & Teardown
	controller, catCtrl, _ := setupTestBoxController("test_boxController.db")
	defer tearDownTestDBController(controller.db)
	boxA, boxB := insertBoxTestData(controller, catCtrl)
	sub := NewMockSubscriber([]string{fmt.Sprintf("box-%d", boxB.ID)})

	// Manipulate boxB
	boxB.Name = "New name for box"
	boxB.Description = "Shiny new description"
	boxB.CategoryID = boxA.CategoryID
	boxB.CreatedAt = time.Now()

	// Update
	err := controller.UpdateBox(boxB.ID, boxB)
	if err != nil {
		t.Log("UpdateBox returned an error", err)
		t.Fail()
	}

	// Load
	box, err := controller.LoadBox(boxB.ID)
	if err != nil {
		t.Log("LoadBox returned an error", err)
		t.Fail()
	}

	// Compare
	if box == nil || !box.Equal(boxB) {
		t.Log("Retrieved box does not equal update", box, boxB)
		t.Fail()
	}

	// Notifications
	if !sub.Assert(fmt.Sprintf("box-%d", boxB.ID), 1) {
		t.Log("Update notification was not sent.")
		t.Fail()
	}

}

func TestBoxCtrlAddBox(t *testing.T) {

	// Setup & Teardown
	controller, catCtrl, _ := setupTestBoxController("test_boxController.db")
	defer tearDownTestDBController(controller.db)
	boxA, _ := insertBoxTestData(controller, catCtrl)
	sub := NewMockSubscriber([]string{"boxes", "stats"})

	// Create boxC
	boxC := models.NewBox()
	boxC.ID = 0
	boxC.Name = "New C-Generation Box"
	boxC.Name = "My Description"
	boxC.CategoryID = boxA.CategoryID

	// Insert
	boxC, err := controller.AddBox(boxC)
	if err != nil {
		t.Log("AddBox returned an error", err)
		t.Fail()
	}

	// Check ID
	if boxC.ID <= 0 {
		t.Log("ID field of inserted box was not updated", boxC.ID)
		t.Fail()
	}

	// Load
	retrievedBox, err := controller.LoadBox(boxC.ID)
	if err != nil {
		t.Log("LoadBox returned an error", err)
		t.Fail()
	}

	// Compare
	if retrievedBox == nil || !retrievedBox.Equal(boxC) {
		t.Log("Retrieved box does not equal inserted", retrievedBox, boxC)
		t.Fail()
	}

	// Notifications
	if !sub.Assert("boxes", 1) {
		t.Log("Insert notification boxes was not sent.")
		t.Fail()
	}
	if !sub.Assert("stats", 1) {
		t.Log("Insert notification stats was not sent.")
		t.Fail()
	}

}

func TestBoxCtrlDeleteBox(t *testing.T) {

	// Setup & Teardown
	controller, catCtrl, qCtrl := setupTestBoxController("test_boxController.db")
	defer tearDownTestDBController(controller.db)
	_, boxB := insertBoxTestData(controller, catCtrl)

	// Insert Questions to BoxB
	insertBoxTestDataQuestionsForBox(boxB.ID, qCtrl)

	// Notifications (after Question generation)
	sub := NewMockSubscriber([]string{"boxes", "stats", "questions"})

	// Delete
	err := controller.DeleteBox(boxB.ID)
	if err != nil {
		t.Log("DeleteBox returned an error", err)
		t.Fail()
	}

	// Check if deleted box is still there
	box, err := controller.LoadBox(boxB.ID)
	if box != nil || err == nil {
		t.Log("Box was not deleted", box)
		t.Fail()
	}

	// Check if there are still questions of that box
	sqlStmt := "SELECT COUNT(*) FROM Question WHERE BoxID = ?;"
	row := controller.db.Connection().QueryRow(sqlStmt, boxB.ID)
	var count int
	err = row.Scan(&count)
	if err != nil {
		t.Log("Error while counting questions", err)
		t.Fail()
	}
	if count != 0 {
		t.Log("Questions of box were not deleted", count)
		t.Fail()
	}

	// Notifications
	if !sub.Assert("boxes", 1) {
		t.Log("Delete notification boxes was not sent.")
		t.Fail()
	}
	if !sub.Assert("stats", 1) {
		t.Log("Delete notification stats was not sent.")
		t.Fail()
	}
	if !sub.Assert("questions", 1) {
		t.Log("Delete notification questions was not sent.")
		t.Fail()
	}

}

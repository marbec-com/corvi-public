package controllers

import (
	"fmt"
	"log"
	"marb.ec/corvi-backend/models"
	"testing"
	"time"
)

func setupTestBoxController(path string, settings *models.Settings) (*BoxController, *CategoryController, *QuestionController) {

	db := setupTestDBController(path)
	s := NewMockSettingsService(settings)

	// Make sure category tables are created
	c, err := NewCategoryController(db, s)
	if err != nil {
		log.Fatal("Error in Setup", err)
		return nil, nil, nil
	}

	// Make sure question tables are created
	q, err := NewQuestionController(db, s)
	if err != nil {
		log.Fatal("Error in Setup", err)
		return nil, nil, nil
	}

	b, err := NewBoxController(db, s)
	if err != nil {
		log.Fatal("Error in Setup", err)
		return nil, nil, nil
	}
	return b, c, q

}

func insertBoxTestData(db DatabaseService) (*models.Box, *models.Box) {

	catA, catB := insertCategoryTestData(db)

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

	insertRawBoxes([]*models.Box{boxA, boxB}, db)

	return boxA, boxB

}

func insertBoxTestDataQuestionsForBox(boxID uint, db DatabaseService) (*models.Question, *models.Question) {

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

	insertRawQuestions([]*models.Question{questionA, questionB}, db)

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
	controller, _, _ := setupTestBoxController("test_boxController.db", nil)
	defer tearDownTestDBController(controller.db)
	boxA, boxB := insertBoxTestData(controller.db)

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
	controller, _, _ := setupTestBoxController("test_boxController.db", nil)
	defer tearDownTestDBController(controller.db)
	_, boxB := insertBoxTestData(controller.db)

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
	controller, _, _ := setupTestBoxController("test_boxController.db", nil)
	defer tearDownTestDBController(controller.db)
	_, boxB := insertBoxTestData(controller.db)

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
	controller, _, _ := setupTestBoxController("test_boxController.db", nil)
	defer tearDownTestDBController(controller.db)
	boxA, boxB := insertBoxTestData(controller.db)
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
	controller, _, _ := setupTestBoxController("test_boxController.db", nil)
	defer tearDownTestDBController(controller.db)
	boxA, _ := insertBoxTestData(controller.db)
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
	controller, _, _ := setupTestBoxController("test_boxController.db", nil)
	defer tearDownTestDBController(controller.db)
	_, boxB := insertBoxTestData(controller.db)

	// Insert Questions to BoxB
	insertBoxTestDataQuestionsForBox(boxB.ID, controller.db)

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

func TestBoxRemoveQuestionFromHeap(t *testing.T) {

	// Setup & Teardown
	controller, _, _ := setupTestBoxController("test_boxController.db", nil)
	defer tearDownTestDBController(controller.db)

	// Create questions
	qA := models.NewQuestion()
	qA.ID = 16
	qA.BoxID = 12
	qB := models.NewQuestion()
	qB.ID = 12
	qB.BoxID = 12
	qC := models.NewQuestion()
	qC.ID = 5
	qC.BoxID = 12

	// Create heaps
	testHeap := models.NewQuestionHeap()
	testHeap.Add(qA)
	testHeap.Add(qB)

	resultHeap := models.NewQuestionHeap()
	resultHeap.Add(qB)

	// Create BoxController and HeapCache
	heapCache := make(map[uint]*models.QuestionHeap, 1)
	heapCache[qA.BoxID] = testHeap
	controller.heapCache = heapCache

	// Perform Remove Question
	err := controller.removeQuestionFromHeap(qA.BoxID, qC.ID)
	if err != nil {
		t.Log("removeQuestionFromHeap returned an error", err)
		t.Fail()
	}
	controller.removeQuestionFromHeap(qA.BoxID, qA.ID)

	// Compare
	if !testHeap.Equal(resultHeap) {
		t.Log("Resulting heap differs from the expected heap", testHeap, resultHeap)
		t.Fail()
	}

}

func TestBoxReAddQuestionFromHeap(t *testing.T) {

	// Setup & Teardown
	controller, _, _ := setupTestBoxController("test_boxController.db", nil)
	defer tearDownTestDBController(controller.db)

	// Create questions
	qA := models.NewQuestion()
	qA.ID = 16
	qA.BoxID = 12
	qB := models.NewQuestion()
	qB.ID = 12
	qB.BoxID = 12
	qC := models.NewQuestion()
	qC.ID = 5
	qC.BoxID = 12

	// Create heaps
	testHeap := models.NewQuestionHeap()
	testHeap.Add(qA)
	testHeap.Add(qB)

	resultHeap := models.NewQuestionHeap()
	resultHeap.Add(qB)
	resultHeap.Add(qA)

	// Create heapCache and assign to BoxController
	heapCache := make(map[uint]*models.QuestionHeap, 1)
	heapCache[qA.BoxID] = testHeap
	controller.heapCache = heapCache

	// Perform Remove Question
	err := controller.reAddQuestionFromHeap(qA.BoxID, qC.ID)
	if err != nil {
		t.Log("reAddQuestionFromHeap returned an error", err)
		t.Fail()
	}
	controller.reAddQuestionFromHeap(qA.BoxID, qA.ID)

	// Compare
	if !testHeap.Equal(resultHeap) {
		t.Log("Resulting heap differs from the expected heap", testHeap, resultHeap)
		t.Fail()
	}

}

func TestBoxGetQuestionToLearn(t *testing.T) {

	// Setup & Teardown
	controller, _, _ := setupTestBoxController("test_boxController.db", nil)
	defer tearDownTestDBController(controller.db)

	// Create questions
	qA := models.NewQuestion()
	qA.ID = 16
	qA.BoxID = 12
	qB := models.NewQuestion()
	qB.ID = 12
	qB.BoxID = 12

	// Create heap
	testHeap := models.NewQuestionHeap()
	testHeap.Add(qA)
	testHeap.Add(qB)

	resultHeap := models.NewQuestionHeap()
	resultHeap.Add(qA)
	resultHeap.Add(qB)

	// Create heapCache and assign to BoxController
	heapCache := make(map[uint]*models.QuestionHeap, 1)
	heapCache[qA.BoxID] = testHeap
	controller.heapCache = heapCache

	// Get top question
	retrievedQuestion, err := controller.GetQuestionToLearn(qA.BoxID)
	if err != nil {
		t.Log("GetQuestionToLearn returned an error", err)
		t.Fail()
	}

	// Test that we got top question
	if !retrievedQuestion.Equal(qA) {
		t.Log("Retrieved question is not top question in heap", retrievedQuestion, testHeap)
		t.Fail()
	}

	// Test that question was not removed from heap
	if !testHeap.Equal(resultHeap) {
		t.Log("Heap was modified", testHeap, resultHeap)
		t.Fail()
	}

}

func TestBoxBuildHeap(t *testing.T) {

	// Setup & Teardown
	settings := models.NewSettings()
	settings.MaxDailyQuestionsPerBox = 3
	settings.RelearnUntilAccomplished = false
	controller, _, _ := setupTestBoxController("test_boxController.db", settings)
	defer tearDownTestDBController(controller.db)
	_, boxB := insertBoxTestData(controller.db)

	// Create questions
	qA := models.NewQuestion()
	qA.BoxID = boxB.ID
	qA.Next = time.Now().AddDate(0, 0, -1)
	qB := models.NewQuestion()
	qB.BoxID = boxB.ID
	qB.Next = time.Now().AddDate(0, 0, -2)
	qC := models.NewQuestion()
	qC.BoxID = boxB.ID
	insertRawQuestions([]*models.Question{qA, qB, qC}, controller.db)

	// Create LearnUnits
	lA := models.NewLearnUnit()
	lA.QuestionID = qC.ID
	lA.BoxID = qC.BoxID
	lA.Correct = true
	lA.PrevCorrect = false
	lB := models.NewLearnUnit()
	lB.QuestionID = qC.ID
	lB.BoxID = qC.BoxID
	lB.Correct = false
	lB.PrevCorrect = false
	insertRawLearnUnits([]*models.LearnUnit{lA, lB}, controller.db)

	// Build Heap
	err := controller.BuildHeap(boxB.ID)
	if err != nil {
		t.Log("buildHeap returned an error", err)
		t.Fail()
	}

	// Compare if new heap in heapCache = expected heap
	expectedHeap := models.NewQuestionHeap()
	expectedHeap.Add(qB)

	// Get result heap
	resultHeap, ok := controller.heapCache[boxB.ID]
	if !ok {
		t.Log("No heap was created", controller.heapCache)
		t.Fail()
	}

	// Compare
	if !expectedHeap.Equal(resultHeap) {
		t.Log("Result heap differs from expected heap", expectedHeap, resultHeap)
		t.Fail()
	}

}

func TestBoxBuildHeapRelearn(t *testing.T) {

	// Setup & Teardown
	settings := models.NewSettings()
	settings.MaxDailyQuestionsPerBox = 3
	settings.RelearnUntilAccomplished = true
	controller, _, _ := setupTestBoxController("test_boxController.db", settings)
	defer tearDownTestDBController(controller.db)
	_, boxB := insertBoxTestData(controller.db)

	// Create questions
	qA := models.NewQuestion()
	qA.BoxID = boxB.ID
	qA.Next = time.Now().AddDate(0, 0, -1)
	qB := models.NewQuestion()
	qB.BoxID = boxB.ID
	qB.Next = time.Now().AddDate(0, 0, -2)
	qC := models.NewQuestion()
	qC.BoxID = boxB.ID
	qD := models.NewQuestion()
	qD.BoxID = boxB.ID
	qD.Next = time.Now().AddDate(0, 0, -3)
	insertRawQuestions([]*models.Question{qA, qB, qC, qD}, controller.db)

	// Create LearnUnits
	lA := models.NewLearnUnit()
	lA.QuestionID = qC.ID
	lA.BoxID = qC.BoxID
	lA.Correct = true
	lA.PrevCorrect = false
	lB := models.NewLearnUnit()
	lB.QuestionID = qC.ID
	lB.BoxID = qC.BoxID
	lB.Correct = false
	lB.PrevCorrect = false
	insertRawLearnUnits([]*models.LearnUnit{lA, lB}, controller.db)

	// Build Heap
	err := controller.BuildHeap(boxB.ID)
	if err != nil {
		t.Log("buildHeap returned an error", err)
		t.Fail()
	}

	// Compare if new heap in heapCache = expected heap
	expectedHeap := models.NewQuestionHeap()
	expectedHeap.Add(qD)
	expectedHeap.Add(qB)

	// Get result heap
	resultHeap, ok := controller.heapCache[boxB.ID]
	if !ok {
		t.Log("No heap was created", controller.heapCache)
		t.Fail()
	}

	// Compare
	if !expectedHeap.Equal(resultHeap) {
		t.Log("Result heap differs from expected heap", expectedHeap, resultHeap)
		t.Fail()
	}

}

func TestBoxBuildHeaps(t *testing.T) {

	// Setup & Teardown
	settings := models.NewSettings()
	settings.MaxDailyQuestionsPerBox = 3
	settings.RelearnUntilAccomplished = false
	controller, _, _ := setupTestBoxController("test_boxController.db", settings)
	defer tearDownTestDBController(controller.db)
	boxA, boxB := insertBoxTestData(controller.db)

	// Create questions
	qA := models.NewQuestion()
	qA.BoxID = boxB.ID
	qA.Next = time.Now().AddDate(0, 0, -1)
	qB := models.NewQuestion()
	qB.BoxID = boxB.ID
	qB.Next = time.Now().AddDate(0, 0, -2)
	qC := models.NewQuestion()
	qC.BoxID = boxB.ID
	insertRawQuestions([]*models.Question{qA, qB, qC}, controller.db)

	// Create LearnUnits
	lA := models.NewLearnUnit()
	lA.QuestionID = qC.ID
	lA.BoxID = qC.BoxID
	lA.Correct = true
	lA.PrevCorrect = false
	lB := models.NewLearnUnit()
	lB.QuestionID = qC.ID
	lB.BoxID = qC.BoxID
	lB.Correct = false
	lB.PrevCorrect = false
	insertRawLearnUnits([]*models.LearnUnit{lA, lB}, controller.db)

	// Build Heap
	err := controller.BuildHeaps()
	if err != nil {
		t.Log("buildHeap returned an error", err)
		t.Fail()
	}

	// Compare if new heap in heapCache = expected heap
	expectedHeap := models.NewQuestionHeap()
	expectedHeap.Add(qB)

	// Get result heaps
	resultHeapB, ok := controller.heapCache[boxB.ID]
	if !ok {
		t.Log("No heapB was created", controller.heapCache)
		t.Fail()
	}
	_, ok = controller.heapCache[boxA.ID]
	if ok {
		t.Log("Heap for BoxA with no questions was created", controller.heapCache)
		t.Fail()
	}

	// Compare
	if !expectedHeap.Equal(resultHeapB) {
		t.Log("Result heap B differs from expected heap", expectedHeap, resultHeapB)
		t.Fail()
	}

}

func TestBoxBuildHeapsRelearn(t *testing.T) {

	// Setup & Teardown
	settings := models.NewSettings()
	settings.MaxDailyQuestionsPerBox = 3
	settings.RelearnUntilAccomplished = true
	controller, _, _ := setupTestBoxController("test_boxController.db", settings)
	defer tearDownTestDBController(controller.db)
	boxA, boxB := insertBoxTestData(controller.db)

	// Create questions
	qA := models.NewQuestion()
	qA.BoxID = boxB.ID
	qA.Next = time.Now().AddDate(0, 0, -1)
	qB := models.NewQuestion()
	qB.BoxID = boxB.ID
	qB.Next = time.Now().AddDate(0, 0, -2)
	qC := models.NewQuestion()
	qC.BoxID = boxB.ID
	qD := models.NewQuestion()
	qD.BoxID = boxB.ID
	qD.Next = time.Now().AddDate(0, 0, -3)
	insertRawQuestions([]*models.Question{qA, qB, qC, qD}, controller.db)

	// Create LearnUnits
	lA := models.NewLearnUnit()
	lA.QuestionID = qC.ID
	lA.BoxID = qC.BoxID
	lA.Correct = true
	lA.PrevCorrect = false
	lB := models.NewLearnUnit()
	lB.QuestionID = qC.ID
	lB.BoxID = qC.BoxID
	lB.Correct = false
	lB.PrevCorrect = false
	insertRawLearnUnits([]*models.LearnUnit{lA, lB}, controller.db)

	// Build Heap
	err := controller.BuildHeaps()
	if err != nil {
		t.Log("buildHeap returned an error", err)
		t.Fail()
	}

	// Compare if new heap in heapCache = expected heap
	expectedHeap := models.NewQuestionHeap()
	expectedHeap.Add(qD)
	expectedHeap.Add(qB)

	// Get result heaps
	resultHeapB, ok := controller.heapCache[boxB.ID]
	if !ok {
		t.Log("No heapB was created", controller.heapCache)
		t.Fail()
	}
	_, ok = controller.heapCache[boxA.ID]
	if ok {
		t.Log("Heap for BoxA with no questions was created", controller.heapCache)
		t.Fail()
	}

	// Compare
	if !expectedHeap.Equal(resultHeapB) {
		t.Log("Result heap differs from expected heap", expectedHeap, resultHeapB)
		t.Fail()
	}

}

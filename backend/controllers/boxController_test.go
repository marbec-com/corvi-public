package controllers

import (
	"log"
	"marb.ec/corvi-backend/models"
	"testing"
	"time"
)

func setupTestBoxController(path string) *BoxController {

	db := setupTestDBController(path)

	// Make sure category tables are created
	_, err := NewCategoryController(db)
	if err != nil {
		log.Fatal("Error in Setup", err)
		return nil
	}

	// Make sure question tables are created
	_, err = NewQuestionController(db)
	if err != nil {
		log.Fatal("Error in Setup", err)
		return nil
	}

	c, err := NewBoxController(db)
	if err != nil {
		log.Fatal("Error in Setup", err)
		return nil
	}
	return c

}

func insertBoxTestData(db *DBController) (*models.Box, *models.Box) {

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

	// SQL INSERT two boxes
	sqlStmt := "INSERT INTO Box (ID, Name, Description, CategoryID, CreatedAt) VALUES (?,?,?,?,?), (?,?,?,?,?);"
	_, err := db.Connection().Exec(sqlStmt, boxA.ID, boxA.Name, boxA.Description, boxA.CategoryID, boxA.CreatedAt, boxB.ID, boxB.Name, boxB.Description, boxB.CategoryID, boxB.CreatedAt)
	if err != nil {
		log.Fatal("Could not insert test data", err)
		return nil, nil
	}

	return boxA, boxB

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
	controller := setupTestBoxController("test_boxController.db")
	defer tearDownTestDBController(controller.db)
	boxA, boxB := insertBoxTestData(controller.db)

	// Create Heap for single Box
	heap := models.NewQuestionHeap()
	heap.Add(models.NewQuestion())
	heap.Add(models.NewQuestion())
	controller.heapCache[boxA.ID] = heap

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
	for _, b := range boxes {
		if b.ID == boxA.ID {
			if b.Name != boxA.Name || b.Description != boxA.Description || b.CategoryID != boxA.CategoryID || !b.CreatedAt.Equal(boxA.CreatedAt) {
				t.Log("First inserted box does not equal result", boxA, b)
				t.Fail()
			}
			// Check heap size
			if b.QuestionsToLearn != uint(heap.Length()) {
				t.Log("Heapsize was not correctly mapped to QuestionsToLearn", b.QuestionsToLearn)
				t.Fail()
			}
		} else if b.ID == boxB.ID {
			if b.Name != boxB.Name || b.Description != boxB.Description || b.CategoryID != boxB.CategoryID || !b.CreatedAt.Equal(boxB.CreatedAt) {
				t.Log("First inserted box does not equal result", boxB, b)
				t.Fail()
			}
		} else {
			t.Log("Unexpected object", b)
			t.Fail()
		}
	}

}

func TestBoxCtrlLoadBoxesOfCategory(t *testing.T) {

	// Setup & Teardown
	controller := setupTestBoxController("test_boxController.db")
	defer tearDownTestDBController(controller.db)
	_, boxB := insertBoxTestData(controller.db)

	// Load Boxes
	boxes, err := controller.LoadBoxesOfCategory(boxB.CategoryID)
	if err != nil {
		t.Log("LoadBoxes returned an error", err)
		t.Fail()
	}

	if len(boxes) != 1 || boxes[0].ID != boxB.ID {
		t.Log("Returned array does not contain boxB", boxes)
		t.Fail()
	}

}

func TestBoxCtrlLoadBox(t *testing.T) {

}

func TestBoxCtrlUpdateBox(t *testing.T) {

}

func TestBoxCtrlAddBox(t *testing.T) {

}

func TestBoxCtrlDeleteBox(t *testing.T) {

}

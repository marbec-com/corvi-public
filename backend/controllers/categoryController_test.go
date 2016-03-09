package controllers

import (
	"fmt"
	"log"
	"marb.ec/corvi-backend/models"
	"testing"
	"time"
)

func setupTestCategoryController(path string) *CategoryController {
	db := setupTestDBController(path)
	c, err := NewCategoryController(db)
	if err != nil {
		log.Fatal("Error in Setup", err)
		return nil
	}
	return c
}

func insertCategoryTestData(catCtrl *CategoryController) (*models.Category, *models.Category) {
	now := time.Now()
	catA := models.NewCategory()
	catA.ID = 1
	catA.Name = "Category A"
	catA.CreatedAt = now
	catB := models.NewCategory()
	catB.ID = 2
	catB.Name = "Category B"
	catB.CreatedAt = now

	// SQL INSERT two categories
	sqlStmt := "INSERT INTO Category (ID, Name, CreatedAt) VALUES (?,?,?), (?,?,?);"
	_, err := catCtrl.db.Connection().Exec(sqlStmt, catA.ID, catA.Name, catA.CreatedAt, catB.ID, catB.Name, catB.CreatedAt)
	if err != nil {
		log.Fatal("Could not insert test data", err)
		return nil, nil
	}

	return catA, catB
}

func TestCategoryCtrlCreateTables(t *testing.T) {

	// Setup & Teardown
	db := setupTestDBController("test_categoryController.db")
	defer tearDownTestDBController(db)

	// Create CategoryController
	categoryController := &CategoryController{
		db: db,
	}

	// Execute createTables()
	err := categoryController.createTables()
	if err != nil {
		t.Log("Error while executing createTables", err)
		t.Fail()
	}

	// Check SQL
	sqlStmt := "SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'Category'"
	row := categoryController.db.Connection().QueryRow(sqlStmt)
	var count int
	err = row.Scan(&count)
	if err != nil || count != 1 {
		t.Log("Table was not created", err, count)
		t.Fail()
	}

}

func TestCategoryCtrlLoadCategories(t *testing.T) {

	// Setup & Teardown
	controller := setupTestCategoryController("test_categoryController.db")
	defer tearDownTestDBController(controller.db)
	catA, catB := insertCategoryTestData(controller)

	// Load and compare
	categories, err := controller.LoadCategories()
	if err != nil {
		t.Log("LoadCategories returned an error", err)
		t.Fail()
	}

	// Compare
	if len(categories) != 2 {
		t.Log("Returned array does not have length of 2:", len(categories))
		t.Fail()
	}

	if !((categories[0].Equal(catA) && categories[1].Equal(catB)) || categories[0].Equal(catB) && categories[1].Equal(catA)) {
		t.Log("Returned categories do not match", categories, catA, catB)
		t.Fail()
	}

}

func TestCategoryCtrlLoadCategory(t *testing.T) {

	// Setup & Teardown
	controller := setupTestCategoryController("test_categoryController.db")
	defer tearDownTestDBController(controller.db)
	_, catB := insertCategoryTestData(controller)

	// Get second
	cat, err := controller.LoadCategory(catB.ID)
	if err != nil {
		t.Log("LoadCategory returned an error", err)
		t.Fail()
	}

	// Compare
	if cat == nil || !cat.Equal(catB) {
		t.Log("Second inserted category does not equal result", catB, cat)
		t.Fail()
	}

}

func TestCategoryCtrlUpdateCategory(t *testing.T) {

	// Setup & Teardown
	controller := setupTestCategoryController("test_categoryController.db")
	defer tearDownTestDBController(controller.db)
	_, catB := insertCategoryTestData(controller)
	sub := NewMockSubscriber([]string{fmt.Sprintf("category-%d", catB.ID)})

	// Manipulate catB
	catB.Name = "NewName"
	catB.CreatedAt = time.Now()

	// Update
	err := controller.UpdateCategory(catB.ID, catB)
	if err != nil {
		t.Log("UpdateCategory returned an error", err)
		t.Fail()
	}

	// Load
	cat, err := controller.LoadCategory(catB.ID)
	if err != nil {
		t.Log("LoadCategory returned an error", err)
		t.Fail()
	}

	// Compare
	if cat == nil || !cat.Equal(catB) {
		t.Log("Retrieved category does not equal update", cat, catB)
		t.Fail()
	}

	// Notifications
	if !sub.Assert(fmt.Sprintf("category-%d", catB.ID), 1) {
		t.Log("Update notification was not sent.")
		t.Fail()
	}

}

func TestCategoryCtrlAddCategory(t *testing.T) {

	// Setup & Teardown
	controller := setupTestCategoryController("test_categoryController.db")
	defer tearDownTestDBController(controller.db)
	sub := NewMockSubscriber([]string{"categories", "stats"})

	cat := models.NewCategory()
	cat.ID = 0
	cat.Name = "CategoryC"
	cat.CreatedAt = time.Now()

	cat, err := controller.AddCategory(cat)
	if err != nil {
		t.Log("AddCategory returned an error", err)
		t.Fail()
	}

	// Check ID
	if cat.ID <= 0 {
		t.Log("ID field of inserted category was not updated", cat.ID)
		t.Fail()
	}

	// Load
	retrievedCat, err := controller.LoadCategory(cat.ID)
	if err != nil {
		t.Log("LoadCategory returned an error", err)
		t.Fail()
	}

	// Compare
	if retrievedCat == nil || !retrievedCat.Equal(cat) {
		t.Log("Retrieved category does not equal inserted", retrievedCat, cat)
		t.Fail()
	}

	// Notifications
	if !sub.Assert("categories", 1) {
		t.Log("Insert notification categories was not sent.")
		t.Fail()
	}
	if !sub.Assert("stats", 1) {
		t.Log("Insert notification stats was not sent.")
		t.Fail()
	}

}

func TestCategoryCtrlDeleteCategory(t *testing.T) {

	// Setup & Teardown
	controller := setupTestCategoryController("test_categoryController.db")
	defer tearDownTestDBController(controller.db)
	_, catB := insertCategoryTestData(controller)
	sub := NewMockSubscriber([]string{"categories", "stats"})

	// Delete
	err := controller.DeleteCategory(catB.ID)
	if err != nil {
		t.Log("DeleteCategory returned an error", err)
		t.Fail()
	}

	// Check if non deleted category is still there
	retrievedCat, err := controller.LoadCategory(catB.ID)
	if retrievedCat != nil || err == nil {
		t.Log("Category was not deleted", retrievedCat)
		t.Fail()
	}

	// Notifications
	if !sub.Assert("categories", 1) {
		t.Log("Delete notification categories was not sent.")
		t.Fail()
	}
	if !sub.Assert("stats", 1) {
		t.Log("Delete notification stats was not sent.")
		t.Fail()
	}

}

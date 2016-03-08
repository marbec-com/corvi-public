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

func insertCategoryTestData(db *DBController) (*models.Category, *models.Category) {
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
	_, err := db.Connection().Exec(sqlStmt, catA.ID, catA.Name, catA.CreatedAt, catB.ID, catB.Name, catB.CreatedAt)
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
	catA, catB := insertCategoryTestData(controller.db)

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
	for _, c := range categories {
		if c.ID == catA.ID {
			if c.Name != catA.Name || !c.CreatedAt.Equal(catA.CreatedAt) {
				t.Log("First inserted category does not equal result", catA, c)
				t.Fail()
			}
		} else if c.ID == catB.ID {
			if c.Name != catB.Name || !c.CreatedAt.Equal(catB.CreatedAt) {
				t.Log("First inserted category does not equal result", catB, c)
				t.Fail()
			}
		} else {
			t.Log("Unexpected object", c)
			t.Fail()
		}
	}

}

func TestCategoryCtrlLoadCategory(t *testing.T) {

	// Setup & Teardown
	controller := setupTestCategoryController("test_categoryController.db")
	defer tearDownTestDBController(controller.db)
	_, catB := insertCategoryTestData(controller.db)

	// Get second
	cat, err := controller.LoadCategory(catB.ID)
	if err != nil {
		t.Log("LoadCategory returned an error", err)
		t.Fail()
	}

	// Compare
	if cat == nil || cat.ID != catB.ID || cat.Name != catB.Name || !cat.CreatedAt.Equal(catB.CreatedAt) {
		t.Log("Second inserted category does not equal result", catB, cat)
	}

}

func TestCategoryCtrlUpdateCategory(t *testing.T) {

	// Setup & Teardown
	controller := setupTestCategoryController("test_categoryController.db")
	defer tearDownTestDBController(controller.db)
	_, catB := insertCategoryTestData(controller.db)
	sub := NewMockSubscriber([]string{fmt.Sprintf("category-%d", catB.ID)})

	newCat := models.NewCategory()
	newCat.ID = catB.ID
	newCat.Name = "NewName"
	newCat.CreatedAt = time.Now()

	// Update
	err := controller.UpdateCategory(catB.ID, newCat)
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
	if cat == nil || cat.ID != newCat.ID || cat.Name != newCat.Name || !cat.CreatedAt.Equal(newCat.CreatedAt) {
		t.Log("Retrieved category does not equal update", cat, newCat)
	}

	// Notifications
	if c, ok := sub.Notifications[fmt.Sprintf("category-%d", catB.ID)]; !ok || c != 1 {
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

	insertedCat, err := controller.AddCategory(cat)
	if err != nil {
		t.Log("AddCategory returned an error", err)
		t.Fail()
	}

	// Check ID
	if insertedCat.ID == 0 {
		t.Log("ID field of inserted category was not updated", insertedCat.ID)
		t.Fail()
	}

	// Load
	retrievedCat, err := controller.LoadCategory(insertedCat.ID)
	if err != nil {
		t.Log("LoadCategory returned an error", err)
		t.Fail()
	}

	// Compare
	if retrievedCat == nil || retrievedCat.ID != insertedCat.ID || retrievedCat.Name != cat.Name || !retrievedCat.CreatedAt.Equal(cat.CreatedAt) {
		t.Log("Retrieved category does not equal inserted", retrievedCat, insertedCat)
	}

	// Notifications
	if c, ok := sub.Notifications["categories"]; !ok || c != 1 {
		t.Log("Insert notification categories was not sent.")
		t.Fail()
	}
	if c, ok := sub.Notifications["stats"]; !ok || c != 1 {
		t.Log("Insert notification stats was not sent.")
		t.Fail()
	}

}

func TestCategoryCtrlDeleteCategory(t *testing.T) {

	// Setup & Teardown
	controller := setupTestCategoryController("test_categoryController.db")
	defer tearDownTestDBController(controller.db)
	_, catB := insertCategoryTestData(controller.db)
	sub := NewMockSubscriber([]string{"categories", "stats"})

	// Delete
	err := controller.DeleteCategory(catB.ID)
	if err != nil {
		t.Log("DeleteCategory returned an error", err)
		t.Fail()
	}

	// Check if non deleted category is still there
	retrievedCat, _ := controller.LoadCategory(catB.ID)
	if retrievedCat != nil {
		t.Log("Category was not deleted", retrievedCat)
		t.Fail()
	}

	// Notifications
	if c, ok := sub.Notifications["categories"]; !ok || c != 1 {
		t.Log("Delete notification categories was not sent.")
		t.Fail()
	}
	if c, ok := sub.Notifications["stats"]; !ok || c != 1 {
		t.Log("Delete notification stats was not sent.")
		t.Fail()
	}

}

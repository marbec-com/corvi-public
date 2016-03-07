package controllers

import (
	"errors"
	"fmt"
	"marb.ec/corvi-backend/models"
	"marb.ec/maf/events"
)

var mockCategories = []*models.Category{
	&models.Category{
		ID:   1,
		Name: "Computer Science",
	},
	&models.Category{
		ID:   2,
		Name: "Vocabulary",
	},
}
var mockCategoriesID uint = 3

var CategoryControllerSingleton *CategoryController

func CategoryCtrl() *CategoryController {
	return CategoryControllerSingleton
}

type CategoryController struct {
	db *DBController
}

func NewCategoryController(db *DBController) (*CategoryController, error) {
	c := &CategoryController{
		db: db,
	}
	err := c.createTables()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *CategoryController) createTables() error {
	sql := "CREATE TABLE IF NOT EXISTS Category (ID INTEGER PRIMARY KEY ASC NOT NULL, Name VARCHAR (255) NOT NULL, CreatedAt DATETIME NOT NULL);"
	_, err := c.db.Connection().Exec(sql)
	return err
}

func (c *CategoryController) LoadCategories() ([]*models.Category, error) {

	// Select all categories
	sql := "SELECT ID, Name, CreatedAt FROM Category;"
	rows, err := c.db.Connection().Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create empty result set
	var result []*models.Category

	for rows.Next() {
		// Create new Category object
		newCat := &models.Category{}
		// Populate
		err = rows.Scan(&newCat.ID, &newCat.Name, &newCat.CreatedAt)
		if err != nil {
			return nil, err
		}
		// Append to result set
		result = append(result, newCat)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *CategoryController) LoadCategory(id uint) (*models.Category, error) {

	// Select category with matching ID
	sql := "SELECT ID, Name, CreatedAt FROM Category WHERE ID = ?;"
	row := c.db.Connection().QueryRow(sql, id)

	// Create new Category object
	newCat := &models.Category{}

	// Populate
	err := row.Scan(&newCat.ID, &newCat.Name, &newCat.CreatedAt)
	if err != nil {
		return nil, err
	}

	return newCat, nil

}

func (c *CategoryController) LoadBoxes(id uint) ([]*models.Box, error) {

	_, err := c.LoadCategory(id)
	if err != nil {
		return nil, err
	}

	boxes := []*models.Box{}
	for _, box := range mockBoxes {
		if box.CategoryID == id {
			boxes = append(boxes, box)
		}
	}

	return boxes, nil
}

func (c *CategoryController) UpdateCategory(catID uint, cat *models.Category) error {
	// Find category to update
	for k, c := range mockCategories {
		if c.ID == catID {
			// Upate category
			mockCategories[k] = cat
			// Publish event to force client refresh
			events.Events().Publish(events.Topic(fmt.Sprintf("category-%d", catID)), c)
			return nil
		}
	}

	// Return error if category was not found
	return errors.New("Category to update was not found.")
}

func (c *CategoryController) AddCategory(cat *models.Category) (*models.Category, error) {

	// Begin Transaction
	tx, err := c.db.Connection().Begin()
	if err != nil {
		return nil, err
	}

	// Rollback in case of an error
	defer tx.Rollback()

	// Execute insert statement
	sql := "INSERT INTO Category (Name, CreatedAt) VALUES (?, ?);"
	res, err := tx.Exec(sql, cat.Name, cat.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Update objects ID
	newID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	cat.ID = uint(newID)

	// Commit
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Publish event to force client refresh
	events.Events().Publish(events.Topic("categories"), c)
	events.Events().Publish(events.Topic("stats"), c)

	// Return inserted object
	return cat, nil

}

func (c *CategoryController) DeleteCategory(catID uint) error {
	for _, box := range mockBoxes {
		if box.CategoryID == catID {
			return errors.New("Cannot delete category that still has boxes assigned.")
		}
	}

	for k, c := range mockCategories {
		if c.ID == catID {
			mockCategories, mockCategories[len(mockCategories)-1] = append(mockCategories[:k], mockCategories[k+1:]...), nil

			// Publish event to force client refresh
			// We don't need to refresh the boxes, sind there shouldn't be any boxes in that category
			events.Events().Publish(events.Topic("categories"), c)
			events.Events().Publish(events.Topic("stats"), c)

			return nil
		}
	}

	return errors.New("Category not found.")
}

package controllers

import (
	"errors"
	"fmt"
	"marb.ec/corvi-backend/models"
	"marb.ec/maf/events"
)

type CategoryController interface {
	CreateTables() error
	LoadCategories() ([]*models.Category, error)
	LoadCategory(id uint) (*models.Category, error)
	UpdateCategory(catID uint, cat *models.Category) error
	AddCategory(cat *models.Category) (*models.Category, error)
	DeleteCategory(catID uint) error
}

type CategoryControllerImpl struct {
	DatabaseService DatabaseService `inject:""`
	SettingsService SettingsService `inject:""`
}

func NewCategoryController() *CategoryControllerImpl {
	c := &CategoryControllerImpl{}
	return c
}

func (c *CategoryControllerImpl) CreateTables() error {

	// Create table, only if it not already exists
	sql := "CREATE TABLE IF NOT EXISTS Category (ID INTEGER PRIMARY KEY ASC NOT NULL, Name VARCHAR (255) NOT NULL, CreatedAt DATETIME NOT NULL);"
	_, err := c.DatabaseService.Connection().Exec(sql)

	return err

}

func (c *CategoryControllerImpl) LoadCategories() ([]*models.Category, error) {

	// Select all categories
	sql := "SELECT ID, Name, CreatedAt FROM Category;"
	rows, err := c.DatabaseService.Connection().Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create empty result set
	result := make([]*models.Category, 0)

	for rows.Next() {
		// Create new Category object
		newCat := models.NewCategory()
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

func (c *CategoryControllerImpl) LoadCategory(id uint) (*models.Category, error) {

	// Select category with matching ID
	sql := "SELECT ID, Name, CreatedAt FROM Category WHERE ID = ?;"
	row := c.DatabaseService.Connection().QueryRow(sql, id)

	// Create new Category object
	newCat := models.NewCategory()

	// Populate
	err := row.Scan(&newCat.ID, &newCat.Name, &newCat.CreatedAt)
	if err != nil {
		return nil, err
	}

	return newCat, nil

}

func (c *CategoryControllerImpl) UpdateCategory(catID uint, cat *models.Category) error {

	// Begin Transaction
	tx, err := c.DatabaseService.Connection().Begin()
	if err != nil {
		return err
	}

	// Rollback in case of an error
	defer tx.Rollback()

	// Update category
	sql := "UPDATE Category SET Name = ?, CreatedAt = ? WHERE ID = ?;"
	res, err := tx.Exec(sql, cat.Name, cat.CreatedAt, catID)
	if err != nil {
		return err
	}

	// Check if update was performed
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// Return error if no object was updated
	if rows == 0 {
		return errors.New("Category to update was not found.")
	}

	// Commit
	err = tx.Commit()
	if err != nil {
		return err
	}

	// Publish event to force client refresh
	events.Events().Publish(events.Topic(fmt.Sprintf("category-%d", catID)), c)

	return nil

}

func (c *CategoryControllerImpl) AddCategory(cat *models.Category) (*models.Category, error) {

	// Begin Transaction
	tx, err := c.DatabaseService.Connection().Begin()
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

func (c *CategoryControllerImpl) DeleteCategory(catID uint) error {

	// Begin Transaction
	tx, err := c.DatabaseService.Connection().Begin()
	if err != nil {
		return err
	}

	// Rollback in case of an error
	defer tx.Rollback()

	// Execute delete statement
	sql := "DELETE FROM Category WHERE ID = ?;"
	res, err := tx.Exec(sql, catID)
	if err != nil {
		return err
	}

	// Check if delete was performed
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// Return error if no object was deleted
	if rows <= 0 {
		return errors.New("Category could not be deleted.")
	}

	// Commit
	err = tx.Commit()
	if err != nil {
		return err
	}

	// Publish event to force client refresh
	// We don't need to refresh the boxes, since there shouldn't be any boxes in that category
	events.Events().Publish(events.Topic("categories"), c)
	events.Events().Publish(events.Topic("stats"), c)

	return nil

}

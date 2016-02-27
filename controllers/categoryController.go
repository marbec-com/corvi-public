package controllers

import (
	"errors"
	"fmt"
	"marb.ec/corvi-backend/models"
	"marb.ec/maf/events"
)

var mockCategories = []*models.Category{
	&models.Category{1, "Computer Science"},
	&models.Category{2, "Vocabulary"},
}
var mockCategoriesID uint = 3

var CategoryControllerSingleton *CategoryController

type CategoryController struct {
}

func CategoryControllerInstance() *CategoryController {
	if CategoryControllerSingleton == nil {
		CategoryControllerSingleton = NewCategoryController()
	}
	return CategoryControllerSingleton
}

func NewCategoryController() *CategoryController {
	return &CategoryController{}
}

func (c *CategoryController) LoadCategories() ([]*models.Category, error) {
	return mockCategories, nil
}

func (c *CategoryController) LoadCategory(id uint) (*models.Category, error) {
	for _, cat := range mockCategories {
		if id == cat.ID {
			return cat, nil
		}
	}

	return nil, errors.New("Category not found.")
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

func (c *CategoryController) AddCategory(cat *models.Category) error {
	// Temporary - not thread safe!
	cat.ID = mockCategoriesID
	mockCategoriesID++

	// Add category
	mockCategories = append(mockCategories, cat)

	// Publish event to force client refresh
	events.Events().Publish(events.Topic("categories"), c)

	return nil
}

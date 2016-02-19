package controllers

import (
	"errors"
	"marb.ec/corvi-backend/models"
)

var mockCategories = []*models.Category{
	&models.Category{1, "Computer Science"},
	&models.Category{2, "Vocabulary"},
}

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
		if box.Category.ID == id {
			boxes = append(boxes, box)
		}
	}

	return boxes, nil
}

func (c *CategoryController) UpdateCategory(cat *models.Category) error {
	return nil
}

func (c *CategoryController) InsertCategory(cat *models.Category) error {
	return nil
}

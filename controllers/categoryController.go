package controllers

import (
	"marb.ec/corvi-backend/models"
)

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

	return []*models.Category{}, nil
}

func (c *CategoryController) LoadCategory(id uint) (*models.Category, error) {

	return &models.Category{}, nil
}

func (c *CategoryController) LoadBoxes(id uint) ([]*models.Box, error) {

	return []*models.Box{}, nil
}

func (c *CategoryController) UpdateCategory(cat *models.Category) error {
	return nil
}

func (c *CategoryController) InsertCategory(cat *models.Category) error {
	return nil
}

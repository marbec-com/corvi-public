package controllers

import (
	"golang.org/x/net/context"
	"marb.ec/corvi-backend/models"
	"marb.ec/maf/interfaces"
	"net/http"
)

type CategoryController struct {
	Category *models.Category
}

func (c *CategoryController) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	rw.Write([]byte("OK"))
}

func (c *CategoryController) LoadCategories() ([]*models.Category, error) {

	return nil, nil
}

func (c *CategoryController) LoadCategory(id uint) (*models.Category, error) {

	return nil, nil
}

func (c *CategoryController) LoadBoxes(id uint) (*[]models.Box, error) {

	return nil, nil
}

package views

import (
	"encoding/json"
	"golang.org/x/net/context"
	"marb.ec/corvi-backend/controllers"
	"marb.ec/corvi-backend/models"
	"marb.ec/maf/interfaces"
	"net/http"
	"strconv"
)

type CategoriesView struct {
	CategoryController controllers.CategoryController `inject:""`
}

func (v *CategoriesView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	categories, err := v.CategoryController.LoadCategories()

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	enc.Encode(categories)

	if n != nil {
		n(rw, r, ctx)
	}
}

type CategoryView struct {
	CategoryController controllers.CategoryController `inject:""`
}

func (v *CategoryView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Load category by ID
	category, err := v.CategoryController.LoadCategory(uint(id))

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Generate JSON response
	rw.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	enc.Encode(category)

	if n != nil {
		n(rw, r, ctx)
	}
}

type CategoryBoxesView struct {
	BoxController controllers.BoxController `inject:""`
}

func (v *CategoryBoxesView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Load boxes by category ID
	boxes, err := v.BoxController.LoadBoxesOfCategory(uint(id))

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Generate JSON response
	rw.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	enc.Encode(boxes)

	if n != nil {
		n(rw, r, ctx)
	}
}

type CategoryUpdateView struct {
	CategoryController controllers.CategoryController `inject:""`
}

func (v *CategoryUpdateView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	defer r.Body.Close()

	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Load existing object to update
	cat, err := v.CategoryController.LoadCategory(uint(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&cat)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = v.CategoryController.UpdateCategory(uint(id), cat)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusOK)

}

type CategoryAddView struct {
	CategoryController controllers.CategoryController `inject:""`
}

func (v *CategoryAddView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	defer r.Body.Close()

	// Construct object via JSON
	cat := models.NewCategory()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cat)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	cat, err = v.CategoryController.AddCategory(cat)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	enc.Encode(cat)

}

type CategoryDeleteView struct {
	CategoryController controllers.CategoryController `inject:""`
}

func (v *CategoryDeleteView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	err = v.CategoryController.DeleteCategory(uint(id))

	if err != nil && err.Error() == "FOREIGN KEY constraint failed" {
		http.Error(rw, "Cannot delete a category that has still boxes assigned.", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusOK)

}

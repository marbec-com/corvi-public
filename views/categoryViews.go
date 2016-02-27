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

type CategoriesView struct{}

func (v *CategoriesView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	controller := controllers.CategoryControllerInstance()
	categories, err := controller.LoadCategories()

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

type CategoryView struct{}

func (v *CategoryView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Load category by ID
	controller := controllers.CategoryControllerInstance()
	category, err := controller.LoadCategory(uint(id))

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

type CategoryBoxesView struct{}

func (v *CategoryBoxesView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Load boxes by category ID
	controller := controllers.CategoryControllerInstance()
	boxes, err := controller.LoadBoxes(uint(id))

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

type CategoryUpdateView struct{}

func (v *CategoryUpdateView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	defer r.Body.Close()

	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Construct object via JSON
	cat := &models.Category{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&cat)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	controller := controllers.CategoryControllerInstance()
	err = controller.UpdateCategory(uint(id), cat)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
	}

	rw.WriteHeader(http.StatusOK)

}

type CategoryAddView struct{}

func (v *CategoryAddView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	defer r.Body.Close()

	// Construct object via JSON
	cat := models.NewCategory()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cat)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	controller := controllers.CategoryControllerInstance()
	err = controller.AddCategory(cat)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
	}

	rw.WriteHeader(http.StatusCreated)

}

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

type BoxesView struct {
	BoxController controllers.BoxController `inject:""`
}

func (v *BoxesView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	boxes, err := v.BoxController.LoadBoxes()

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	enc.Encode(boxes)

	if n != nil {
		n(rw, r, ctx)
	}
}

type BoxView struct {
	BoxController controllers.BoxController `inject:""`
}

func (v *BoxView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Load box by ID
	box, err := v.BoxController.LoadBox(uint(id))

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Generate JSON response
	rw.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	enc.Encode(box)

	if n != nil {
		n(rw, r, ctx)
	}
}

type BoxQuestionsView struct {
	QuestionController controllers.QuestionController `inject:""`
}

func (v *BoxQuestionsView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Load questions by box ID
	questions, err := v.QuestionController.LoadQuestionsOfBox(uint(id))

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Generate JSON response
	rw.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	enc.Encode(questions)

	if n != nil {
		n(rw, r, ctx)
	}
}

type BoxGetQuestionToLearnView struct {
	BoxController controllers.BoxController `inject:""`
}

func (v *BoxGetQuestionToLearnView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Load box by ID
	question, err := v.BoxController.GetQuestionToLearn(uint(id))

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	if question == nil {
		http.Error(rw, "No questions left for today.", http.StatusNoContent)
		return
	}

	// Generate JSON response
	rw.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	enc.Encode(question)

	if n != nil {
		n(rw, r, ctx)
	}
}

type BoxUpdateView struct {
	BoxController controllers.BoxController `inject:""`
}

func (v *BoxUpdateView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	defer r.Body.Close()

	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Load existing object to update
	box, err := v.BoxController.LoadBox(uint(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&box)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = v.BoxController.UpdateBox(uint(id), box)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusOK)

}

type BoxAddView struct {
	BoxController controllers.BoxController `inject:""`
}

func (v *BoxAddView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	defer r.Body.Close()

	// Construct object via JSON
	box := models.NewBox()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&box)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	box, err = v.BoxController.AddBox(box)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	enc.Encode(box)

}

type BoxDeleteView struct {
	BoxController controllers.BoxController `inject:""`
}

func (v *BoxDeleteView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	err = v.BoxController.DeleteBox(uint(id))

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

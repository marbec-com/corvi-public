package views

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"marb.ec/corvi-backend/controllers"
	"marb.ec/corvi-backend/models"
	"marb.ec/maf/interfaces"
	"net/http"
	"strconv"
)

type QuestionsView struct {
	QuestionController controllers.QuestionController `inject:""`
}

func (v *QuestionsView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	questions, err := v.QuestionController.LoadQuestions()

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	enc.Encode(questions)

	if n != nil {
		n(rw, r, ctx)
	}
}

type QuestionView struct {
	QuestionController controllers.QuestionController `inject:""`
}

func (v *QuestionView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Load question by ID
	question, err := v.QuestionController.LoadQuestion(uint(id))

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
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

type QuestionGiveCorrectAnswerView struct {
	QuestionController controllers.QuestionController `inject:""`
}

func (v *QuestionGiveCorrectAnswerView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Call method by ID
	err = v.QuestionController.GiveAnswer(uint(id), true)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprint(rw, "true")

	if n != nil {
		n(rw, r, ctx)
	}
}

type QuestionGiveWrongAnswerView struct {
	QuestionController controllers.QuestionController `inject:""`
}

func (v *QuestionGiveWrongAnswerView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Call method by ID
	err = v.QuestionController.GiveAnswer(uint(id), false)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprint(rw, "true")

	if n != nil {
		n(rw, r, ctx)
	}
}

type QuestionUpdateView struct {
	QuestionController controllers.QuestionController `inject:""`
}

func (v *QuestionUpdateView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	defer r.Body.Close()

	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	question := models.NewQuestion()

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&question)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = v.QuestionController.UpdateQuestion(uint(id), question)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusOK)

}

type QuestionAddView struct {
	QuestionController controllers.QuestionController `inject:""`
}

func (v *QuestionAddView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	defer r.Body.Close()

	// Construct object via JSON
	question := models.NewQuestion()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&question)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	question, err = v.QuestionController.AddQuestion(question)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	rw.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	enc.Encode(question)

}

type QuestionDeleteView struct {
	QuestionController controllers.QuestionController `inject:""`
}

func (v *QuestionDeleteView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	err = v.QuestionController.DeleteQuestion(uint(id))

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

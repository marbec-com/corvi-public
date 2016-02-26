package views

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"marb.ec/corvi-backend/controllers"
	"marb.ec/maf/interfaces"
	"net/http"
	"strconv"
)

type QuestionsView struct{}

func (v *QuestionsView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	controller := controllers.QuestionControllerInstance()
	questions, err := controller.LoadQuestions()

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

type QuestionView struct{}

func (v *QuestionView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Load question by ID
	controller := controllers.QuestionControllerInstance()
	question, err := controller.LoadQuestion(uint(id))

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

type QuestionGiveCorrectAnswerView struct{}

func (v *QuestionGiveCorrectAnswerView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Call method by ID
	controller := controllers.QuestionControllerInstance()
	err = controller.GiveCorrectAnswer(uint(id))

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprint(rw, "true")

	if n != nil {
		n(rw, r, ctx)
	}
}

type QuestionGiveWrongAnswerView struct{}

func (v *QuestionGiveWrongAnswerView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	// Parse and convert ID
	idRaw := ctx.Value("id").(string)
	id, err := strconv.ParseUint(idRaw, 10, 32)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	// Call method by ID
	controller := controllers.QuestionControllerInstance()
	err = controller.GiveWrongAnswer(uint(id))

	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprint(rw, "true")

	if n != nil {
		n(rw, r, ctx)
	}
}

package views

import (
	"encoding/json"
	"golang.org/x/net/context"
	"marb.ec/corvi-backend/controllers"
	"marb.ec/maf/interfaces"
	"net/http"
)

type BoxesView struct{}

func (v *BoxesView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	controller := controllers.BoxControllerInstance()
	boxes, err := controller.LoadBoxes()

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

type BoxView struct{}

func (v *BoxView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	rw.Write([]byte("OK"))

	if n != nil {
		n(rw, r, ctx)
	}
}

type BoxQuestionsView struct{}

func (v *BoxQuestionsView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	rw.Write([]byte("OK"))

	if n != nil {
		n(rw, r, ctx)
	}
}

type BoxGetQuestionToLearnView struct{}

func (v *BoxGetQuestionToLearnView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	rw.Write([]byte("OK"))

	if n != nil {
		n(rw, r, ctx)
	}
}

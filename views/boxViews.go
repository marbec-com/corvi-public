package views

import (
	"golang.org/x/net/context"
	"marb.ec/maf/interfaces"
	"net/http"
)

type BoxesView struct{}

func (v *BoxesView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	rw.Write([]byte("OK"))

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

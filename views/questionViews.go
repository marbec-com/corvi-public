package views

import (
	"golang.org/x/net/context"
	"marb.ec/maf/interfaces"
	"net/http"
)

type QuestionsView struct{}

func (v *QuestionsView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	rw.Write([]byte("OK"))

	if n != nil {
		n(rw, r, ctx)
	}
}

type QuestionView struct{}

func (v *QuestionView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	rw.Write([]byte("OK"))

	if n != nil {
		n(rw, r, ctx)
	}
}

type QuestionGiveCorrectAnswerView struct{}

func (v *QuestionGiveCorrectAnswerView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	rw.Write([]byte("OK"))

	if n != nil {
		n(rw, r, ctx)
	}
}

type QuestionGiveWrongAnswerView struct{}

func (v *QuestionGiveWrongAnswerView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	rw.Write([]byte("OK"))

	if n != nil {
		n(rw, r, ctx)
	}
}

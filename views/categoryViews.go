package views

import (
	"golang.org/x/net/context"
	"marb.ec/maf/interfaces"
	"net/http"
)

type CategoriesView struct{}

func (v *CategoriesView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	rw.Write([]byte("OK"))

	if n != nil {
		n(rw, r, ctx)
	}
}

type CategoryView struct{}

func (v *CategoryView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	rw.Write([]byte("OK"))

	if n != nil {
		n(rw, r, ctx)
	}
}

type CategoryBoxesView struct{}

func (v *CategoryBoxesView) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	rw.Write([]byte("OK"))

	if n != nil {
		n(rw, r, ctx)
	}
}

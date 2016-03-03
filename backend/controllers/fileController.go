package controllers

import (
	"golang.org/x/net/context"
	"marb.ec/maf/interfaces"
	"net/http"
	"path"
)

type IndexController struct{}

func (c *IndexController) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	http.ServeFile(rw, r, path.Join("frontend", "index.html"))

	if n != nil {
		n(rw, r, ctx)
	}
}

type FileController struct{}

func (c *FileController) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {

	fileHandler := http.FileServer(http.Dir("frontend"))
	fileHandler = http.StripPrefix("/app/", fileHandler)
	fileHandler.ServeHTTP(rw, r)

	if n != nil {
		n(rw, r, ctx)
	}
}

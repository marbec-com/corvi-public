package middleware

import (
	"golang.org/x/net/context"
	"log"
	"marb.ec/maf/interfaces"
	"net/http"
)

type NotFoundHandler struct{}

func (nf *NotFoundHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	http.NotFound(rw, r)
	log.Printf("NOT FOUND: %s %s %s%s from %s\n", r.Proto, r.Method, r.Host, r.URL, r.RemoteAddr)
	n(rw, r, ctx)
}

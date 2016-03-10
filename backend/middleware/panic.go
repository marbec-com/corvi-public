package middleware

import (
	"golang.org/x/net/context"
	"log"
	"marb.ec/maf/interfaces"
	"net/http"
)

type PanicRecoveryHandler struct{}

func (p *PanicRecoveryHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Printf("Panic %v in request %v", err, r)
		}
	}()
	n(rw, r, ctx)
}

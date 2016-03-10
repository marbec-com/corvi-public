package middleware

import (
	"golang.org/x/net/context"
	"log"
	"marb.ec/maf/interfaces"
	"net/http"
	"time"
)

type LogHandler struct{}

func (l *LogHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx context.Context, n interfaces.HandlerFunc) {
	start := time.Now()
	n(rw, r, ctx)
	end := time.Now()
	log.Printf("%s %s %s%s from %s in %v\n", r.Proto, r.Method, r.Host, r.URL, r.RemoteAddr, end.Sub(start))
}

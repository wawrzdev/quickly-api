package services

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"quickly-api/pkg"
	"time"
)

type LogHTTPMiddleware struct {
	l *log.Logger
}

func NewLogHTTPMiddleware(l *log.Logger) *LogHTTPMiddleware {
	return &LogHTTPMiddleware{l}
}

func (m *LogHTTPMiddleware) Func() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			logRw := pkg.NewLogResponseWriter(w)
			next.ServeHTTP(logRw, r)

			//TODO: current route isn't attached to r since it hasnt gone through mux?
			m.l.Printf("%s Duration: %s Status: %d Body: %s)", mux.CurrentRoute(r).GetName(), time.Since(startTime).String(), logRw.StatusCode, logRw.Buf.String())
		})
	}
}

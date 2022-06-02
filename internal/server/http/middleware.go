package internalhttp

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Bytes      int
}

func loggingMiddleware(next http.Handler, logger Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrt := &ResponseWriter{w, 0, 0}
		next.ServeHTTP(wrt, r)
		logger.LogHTTP(r, wrt.StatusCode, wrt.Bytes)
	})
}

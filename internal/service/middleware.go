package service

import (
	"context"
	log "log/slog"
	"net/http"

	"github.com/google/uuid"
)

func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Context().Value(requestIDKey{}).(string)

		log.Info("Request", "method", r.Method, "url", r.URL.String(), "request_id", rid)

		next.ServeHTTP(w, r)
	})
}

func logResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := newResponseWriter(w)

		next.ServeHTTP(rw, r)

		rid := r.Context().Value(requestIDKey{}).(string)

		log.Info("Response", "status", rw.status, "request_id", rid)
	})
}

type requestIDKey struct{}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get("X-Request-ID")

		if rid == "" {
			// New random request ID.
			rid = newRequestID()

			r.Header.Set("X-Request-ID", rid)
		}

		ctx := r.Context()

		ctx = context.WithValue(ctx, requestIDKey{}, rid)

		w.Header().Set("X-Request-ID", rid)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func newRequestID() string {
	u := uuid.New()

	return u.String()
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status

	rw.ResponseWriter.WriteHeader(status)
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Panic recovered", "error", err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

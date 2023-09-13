package main

import (
	"context"
	"errors"
	log "log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	mux := http.NewServeMux()

	mux.Handle("/", handler())

	port := "8080"

	log.Info("Starting server", "port", port)

	server := &http.Server{
		Addr:                         net.JoinHostPort("", port),
		Handler:                      mux,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  0,
		ReadHeaderTimeout:            0,
		WriteTimeout:                 0,
		IdleTimeout:                  0,
		MaxHeaderBytes:               0,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     nil,
		BaseContext:                  nil,
		ConnContext:                  nil,
	}

	var wg sync.WaitGroup

	wg.Add(1)

	server.RegisterOnShutdown(func() {
		defer wg.Done()
		log.Info("Server shutting down")

		server.SetKeepAlivesEnabled(false)

		log.Info("Server shutdown complete")
	})

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Error starting server", "error", err)
			cancel()
		}
	}()

	<-ctx.Done()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Error shutting down server", "error", err)
	}

	wg.Wait()

	log.Info("Exit", "signal", ctx.Err())
}

func handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("Request received", "method", r.Method, "url", r.URL)
		_, err := w.Write([]byte("Hello World"))
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	})
}

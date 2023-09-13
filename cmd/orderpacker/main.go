package main

import (
	"context"
	"errors"
	"fmt"
	log "log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/obalunenko/orderpacker/internal/packer"
	"github.com/obalunenko/orderpacker/internal/service"
)

func main() {
	signals := make(chan os.Signal, 1)

	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	signal.Notify(signals, os.Interrupt, os.Kill)

	go func() {
		s := <-signals

		cancel(fmt.Errorf("received signal: %s", s.String()))
	}()

	port := "8080"

	r := service.NewRouter(packer.NewPacker(packer.WithDefaultBoxes()))

	log.Info("Starting server", "port", port)

	server := &http.Server{
		Addr:                         net.JoinHostPort("", port),
		Handler:                      r,
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
			cancel(err)
		}
	}()

	<-ctx.Done()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Error shutting down server", "error", err)
	}

	wg.Wait()

	log.Info("Exit", "cause", context.Cause(ctx))
}

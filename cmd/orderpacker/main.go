package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/obalunenko/logger"
	_ "github.com/swaggo/swag"

	"github.com/obalunenko/orderpacker/internal/config"
	"github.com/obalunenko/orderpacker/internal/packer"
	"github.com/obalunenko/orderpacker/internal/service"
)

var errSignal = errors.New("received signal")

// @title						Order Packer API
// @version					1.0
// @description				This is a simple API for packing orders
// @termsOfService				http://swagger.io/terms/
// @contact.name				Oleg Balunenko
// @contact.email				oleg.balunenko@gmail.com
// @license.name				MIT
// @license.url				https://opensource.org/license/mit
// @host						localhost:8080
// @schemes					http
//
// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	signals := make(chan os.Signal, 1)

	l := log.FromContext(context.Background())

	ctx := log.ContextWithLogger(context.Background(), l)

	printVersion(ctx)

	ctx, cancel := context.WithCancelCause(ctx)
	defer func() {
		const msg = "Exit"

		var code int

		err := context.Cause(ctx)
		if err != nil && !errors.Is(err, errSignal) {
			code = 1
		}

		l := log.WithField(ctx, "cause", err)

		if code == 0 {
			l.Info(msg)

			return
		}

		l.Error(msg)

		os.Exit(code)
	}()

	defer cancel(nil)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	go func() {
		s := <-signals

		cancel(fmt.Errorf("%w: %s", errSignal, s.String()))
	}()

	cfg, err := config.Load(ctx)
	if err != nil {
		cancel(fmt.Errorf("failed to load config: %w", err))

		return
	}

	l = log.Init(ctx, log.Params{
		Writer: os.Stderr,
		Level:  cfg.Log.Level,
		Format: cfg.Log.Format,
	})

	ctx = log.ContextWithLogger(ctx, l)

	port := cfg.HTTP.Port
	host := cfg.HTTP.Host

	p, err := packer.NewPacker(ctx, packer.WithBoxes(cfg.Pack.Boxes))
	if err != nil {
		cancel(fmt.Errorf("failed to create packer: %w", err))

		return
	}

	log.WithFields(ctx, log.Fields{
		"host": host,
		"port": port,
	}).Info("Starting server")

	server := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: service.NewRouter(p),
	}

	var wg sync.WaitGroup

	wg.Add(1)

	server.RegisterOnShutdown(func() {
		defer wg.Done()
		log.Info(ctx, "Server shutting down")

		server.SetKeepAlivesEnabled(false)

		log.Info(ctx, "Server shutdown complete")
	})

	go func() {
		log.WithField(ctx, "addr", server.Addr).Info("Server started")

		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			cancel(fmt.Errorf("failed to start server: %w", err))
		}
	}()

	<-ctx.Done()

	if err = server.Shutdown(ctx); err != nil {
		log.WithError(ctx, err).Error("Error shutting down server")
	}

	wg.Wait()
}

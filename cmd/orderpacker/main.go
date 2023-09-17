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

	"github.com/obalunenko/getenv"
	log "github.com/obalunenko/logger"

	"github.com/obalunenko/orderpacker/internal/config"
	"github.com/obalunenko/orderpacker/internal/packer"
	"github.com/obalunenko/orderpacker/internal/service"
)

const (
	configPathEnv = "ORDERPACKER_CONFIG_PATH"
)

var errSignal = errors.New("received signal")

func main() {
	printVersion()

	signals := make(chan os.Signal, 1)

	l := log.FromContext(context.Background())

	ctx := log.ContextWithLogger(context.Background(), l)

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

	signal.Notify(signals, os.Interrupt, os.Kill)

	go func() {
		s := <-signals

		cancel(fmt.Errorf("%w: %s", errSignal, s.String()))
	}()

	var useDefaultConfig bool

	cfgPath, err := getenv.Env[string](configPathEnv)
	if err != nil {
		if errors.Is(err, getenv.ErrNotSet) {
			log.WithField(ctx, "env", configPathEnv).Warn("Config path env not set")

			useDefaultConfig = true
		}
	}

	var cfg *config.Config

	if !useDefaultConfig {
		log.WithField(ctx, "path", cfgPath).Info("Using config")

		cfg, err = config.Load(cfgPath)
		if err != nil {
			cancel(fmt.Errorf("failed to load config: %w", err))

			return
		}
	} else {
		log.Warn(ctx, "Using default config")

		cfg = config.DefaultConfig()
	}

	l = log.Init(ctx, log.Params{
		Writer:     os.Stderr,
		Level:      cfg.Log.Level,
		Format:     cfg.Log.Format,
		WithSource: false,
	})

	ctx = log.ContextWithLogger(ctx, l)

	port := cfg.HTTP.Port

	p, err := packer.NewPacker(ctx, packer.WithBoxes(cfg.Pack.Boxes))
	if err != nil {
		cancel(fmt.Errorf("failed to create packer: %w", err))

		return
	}

	r := service.NewRouter(p)

	log.WithField(ctx, "port", port).Info("Starting server")

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
		log.Info(ctx, "Server shutting down")

		server.SetKeepAlivesEnabled(false)

		log.Info(ctx, "Server shutdown complete")
	})

	go func() {
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

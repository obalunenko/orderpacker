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
	"strings"
	"sync"

	"github.com/obalunenko/getenv"

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

	ctx, cancel := context.WithCancelCause(context.Background())
	defer func() {
		const msg = "Exit"

		var code int

		err := context.Cause(ctx)
		if err != nil && !errors.Is(err, errSignal) {
			code = 1
		}

		l := log.With("cause", err)

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
			log.Warn("Config path env not set", "env", configPathEnv)

			useDefaultConfig = true
		}
	}

	var cfg *config.Config

	if !useDefaultConfig {
		log.Info("Using config", "path", cfgPath)

		cfg, err = config.Load(cfgPath)
		if err != nil {
			cancel(fmt.Errorf("failed to load config: %w", err))

			return
		}
	} else {
		log.Warn("Using default config")

		cfg = config.DefaultConfig()
	}

	setLogger(cfg)

	port := cfg.HTTP.Port

	p, err := packer.NewPacker(packer.WithBoxes(cfg.Pack.Boxes))
	if err != nil {
		cancel(fmt.Errorf("failed to create packer: %w", err))

		return
	}

	r := service.NewRouter(p)

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
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			cancel(fmt.Errorf("failed to start server: %w", err))
		}
	}()

	<-ctx.Done()

	if err = server.Shutdown(ctx); err != nil {
		log.Error("Error shutting down server", "error", err)
	}

	wg.Wait()
}

func setLogger(cfg *config.Config) {
	var level log.Leveler

	switch strings.ToLower(cfg.Log.Level) {
	case "debug":
		level = log.LevelDebug
	case "info":
		level = log.LevelInfo
	case "warn":
		level = log.LevelWarn
	case "error":
		level = log.LevelError
	default:
		log.Warn("Unknown log level, info will be used", "level", cfg.Log.Level)

		level = log.LevelInfo
	}

	var handler log.Handler

	switch strings.ToLower(cfg.Log.Format) {
	case "json":
		handler = log.NewJSONHandler(os.Stdout, &log.HandlerOptions{
			AddSource:   false,
			Level:       level,
			ReplaceAttr: nil,
		})
	case "text":
		handler = log.NewTextHandler(os.Stdout, &log.HandlerOptions{
			AddSource:   false,
			Level:       level,
			ReplaceAttr: nil,
		})
	default:
		log.Warn("Unknown log format, text will be used", "format", cfg.Log.Format)

		handler = log.NewTextHandler(os.Stdout, &log.HandlerOptions{
			AddSource:   false,
			Level:       level,
			ReplaceAttr: nil,
		})
	}

	log.SetDefault(log.New(handler))

	log.Info("Logger set", "level", level, "format", cfg.Log.Format)
}

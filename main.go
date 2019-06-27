package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bcmendoza/xds-explorer/handlers"
	"github.com/bcmendoza/xds-explorer/models"
	"github.com/bcmendoza/xds-explorer/stream"

	"github.com/rs/zerolog"
)

func main() {
	var err error

	// startup
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp}).
		With().Timestamp().Str("service", "xds-explorer").Logger()
	logger.Info().Msg("Startup XDS Explorer")

	// signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	logger.Info().Msg("Watch OS signals")

	// context
	ctx, cancelFunc := context.WithCancel(context.Background())

	// request channel
	requestChan := make(chan models.Request, 1)

	// GRPC stream
	streamLogger := logger.With().Str("package", "stream").Logger()
	go stream.Listen(ctx, requestChan, streamLogger)

	// REST server
	serverLogger := logger.With().Str("package", "handlers").Logger()
	server := http.Server{
		Addr:    "0.0.0.0:3001",
		Handler: handlers.Handlers(requestChan, serverLogger),
	}
	go func() {
		serverLogger.Info().Msg("Startup REST server")
		if err = server.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
			serverLogger.Error().AnErr("server.ListenAndServe()", err).Msg("REST server error")
		}
	}()

	// shutdown
	s := <-sigChan
	cancelFunc()
	if err = server.Close(); err != nil {
		logger.Error().AnErr("server.Close()", err).Msg("REST server shutdown error")
	} else {
		logger.Info().Msg("Shutdown REST server")
	}
	logger.Info().Str("signal", s.String()).Msg("Shutdown XDS Explorer")
}

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bcmendoza/xds-explorer/handlers"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func main() {
	var err error

	// logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp}).
		With().Timestamp().Str("service", "xds-explorer").Logger()
	logger.Info().Msg("Starting XDS Explorer")

	// signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// set environment defaults
	viper.SetDefault("port", 8080)

	// get log level
	if viper.GetBool("debug") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Start server
	server := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", viper.GetInt("port")),
		Handler: handlers.Handlers(logger),
	}
	if err := server.ListenAndServe(); err != nil {
		logger.Error().AnErr("server.ListenAndServe()", err).Msg("Start server")
	}

	s := <-sigChan
	logger.Info().Str("signal", s.String()).Msg("shutting down")
	if err = server.Close(); err != nil {
		logger.Debug().AnErr("server.Close()", err).Msg("shutting down")
	}
}

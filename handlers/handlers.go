package handlers

import (
	"fmt"
	"net/http"

	"github.com/bcmendoza/xds-explorer/models"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

func Handlers(requestChan chan<- models.Request, logger zerolog.Logger) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/request", request(requestChan, logger))
	r.HandleFunc("/listen", listen(logger))
	r.HandleFunc("/ping", ping(logger))
	return r
}

func request(requestChan chan<- models.Request, logger zerolog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: change to POST
		if logger, ok := verifyMethod("/request", r.Method, "GET", logger, w); ok {

			// TODO: validate fields before sending to requestChan
			requestChan <- models.Request{
				ResourceType:  "type.googleapis.com/envoy.api.v2.ClusterLoadAssignment",
				Zone:          "default-zone",
				Cluster:       "catalog",
				ResourceNames: []string{"catalog"},
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			jsonResp := "{\"requested\": true}"
			if _, err := w.Write([]byte(jsonResp)); err != nil {
				logger.Error().AnErr("w.Write", err).Msg("500 Internal server error")
			} else {
				logger.Info().Msg("200 OK")
			}
		}
	}
}

func listen(logger zerolog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if logger, ok := verifyMethod("/listen", r.Method, "GET", logger, w); ok {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			// TODO: return internal model resources as JSON
			jsonResp := "{\"listening\": true}"
			if _, err := w.Write([]byte(jsonResp)); err != nil {
				logger.Error().AnErr("w.Write", err).Msg("500 Internal server error")
			} else {
				logger.Info().Msg("200 OK")
			}
		}
	}
}

func ping(logger zerolog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if logger, ok := verifyMethod("/ping", r.Method, "GET", logger, w); ok {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			jsonResp := "{\"ping\": \"pong\"}"
			if _, err := w.Write([]byte(jsonResp)); err != nil {
				logger.Error().AnErr("w.Write", err).Msg("500 Internal server error")
			} else {
				logger.Info().Msg("200 OK")
			}
		}
	}
}

func verifyMethod(route, method, expectedMethod string, logger zerolog.Logger, w http.ResponseWriter) (zerolog.Logger, bool) {
	logger = logger.With().Str("request-type", fmt.Sprintf("%s:'%s'", method, route)).Logger()
	if method != expectedMethod {
		logger.Warn().Msg("405 Method Not Allowed")
		Report(ProblemDetail{StatusCode: http.StatusMethodNotAllowed, Detail: method}, w)
		return logger, false
	}
	return logger, true
}

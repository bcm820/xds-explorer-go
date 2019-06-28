package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bcmendoza/xds-explorer/model"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

func Handlers(requestChan chan<- model.Request, resources *model.Resources, logger zerolog.Logger) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/request", request(requestChan, logger))
	r.HandleFunc("/listen", listen(resources, logger))
	r.HandleFunc("/ping", ping(logger))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("/app/client")))
	return r
}

func request(requestChan chan<- model.Request, logger zerolog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if logger, ok := verifyMethod("/request", r.Method, "POST", logger, w); ok {
			var request model.Request
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&request)
			if err != nil {
				logger.Error().AnErr("json.NewDecoder", err).Msg("400 Bad Request")
				Report(ProblemDetail{
					StatusCode: http.StatusBadRequest,
					Detail:     "Could not unmarshall request JSON",
				}, w)
				return
			}

			requestChan <- request
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			jsonResp := "{\"request updated\": true}"
			if _, err := w.Write([]byte(jsonResp)); err != nil {
				logger.Error().AnErr("w.Write", err).Msg("500 Internal server error")
			} else {
				logger.Info().Msg("200 OK")
			}
		}
	}
}

func listen(resources *model.Resources, logger zerolog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if logger, ok := verifyMethod("/listen", r.Method, "GET", logger, w); ok {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")

			collection := resources.GetCLAs()
			if collection == nil {
				collection = make([]v2.ClusterLoadAssignment, 0)
			}

			jsonResp, err := json.Marshal(collection)
			if err != nil {
				logger.Error().AnErr("json.Marshal", err).Msg("Could not marshal into JSON")
			}

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
	logger = logger.With().Str("request-type", fmt.Sprintf("%s %s", method, route)).Logger()
	if method != expectedMethod {
		logger.Warn().Msg("405 Method Not Allowed")
		Report(ProblemDetail{StatusCode: http.StatusMethodNotAllowed, Detail: method}, w)
		return logger, false
	}
	logger.Info().Msg("Receive request")
	return logger, true
}

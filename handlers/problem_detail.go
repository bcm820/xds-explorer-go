package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

// ProblemDetail is the Decipher format of the rfc 7807 Problem details
// See https://tools.ietf.org/html/rfc7807
type ProblemDetail struct {

	// Type (optional) is a URI giving detailed defintion of the problem
	Type string `json:"type,omitempty"`

	// StatusCode is the numeric HTTP Status
	StatusCode int `json:"code"`

	// Title is the (short) title of the problem, generally the HTTP status
	// If blank, will be replaced by http.StatusText(StatusCode)
	Title string `json:"title"`

	// Detail (optional) detailed information about the
	Detail string `json:"detail,omitempty"`
}

// Report report sends the problem object as the body of the http.Response
// with the appropriate 'Content-Type: application/problem+json' header
func Report(problem ProblemDetail, w http.ResponseWriter) {
	if problem.Title == "" {
		problem.Title = http.StatusText(problem.StatusCode)
	}
	w.Header().Add("Content-Type", "application/problem+json")
	w.WriteHeader(problem.StatusCode)
	json.NewEncoder(w).Encode(&problem)
}

// Unmarshal unmarshals the body of a http.Response with
// 'Content-Type: application/problem+json'
func Unmarshal(body io.Reader) (ProblemDetail, error) {
	var problem ProblemDetail
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&problem)
	return problem, err
}

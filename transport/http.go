package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/benkim0414/geoauth2/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler")
)

// NewHTTPHandler returns an HTTP handler that makes a set of endpoints
// available on predefined paths.
func NewHTTPHandler(endpoints endpoint.Endpoints, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorLogger(logger),
	}
	r := mux.NewRouter().PathPrefix("/api/v0/").Subrouter()
	r.Methods("POST").Path("/clients").Handler(httptransport.NewServer(
		endpoints.PostClientEndpoint,
		decodeHTTPPostClientRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	r.Methods("GET").Path("/clients/{id}").Handler(httptransport.NewServer(
		endpoints.GetClientEndpoint,
		decodeHTTPGetClientRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	r.Methods("DELETE").Path("/clients/{id}").Handler(httptransport.NewServer(
		endpoints.DeleteClientEndpoint,
		decodeHTTPDeleteClientRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	return r
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func err2code(err error) int {
	return http.StatusInternalServerError
}

func errorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

type errorWrapper struct {
	Error string `json:"error"`
}

// decodeHTTPPostClientRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded POST client request from the HTTP request body. Primarily useful in a
// server.
func decodeHTTPPostClientRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.PostClientRequest
	err := json.NewDecoder(r.Body).Decode(&req.Client)
	return req, err
}

// decodeHTTPGetClientRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded GET client request from the HTTP request body. Primarily useful in a
// server.
func decodeHTTPGetClientRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return endpoint.GetClientRequest{ID: id}, nil
}

// decodeHTTPDeleteClientRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded DELETE client request from the HTTP request body. Primarily useful in a
// server.
func decodeHTTPDeleteClientRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return endpoint.DeleteClientRequest{ID: id}, nil
}

// encodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

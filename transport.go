package scoresvcdemo

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"net/url"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Expect this to be returned only on GET request
// with {:id} url parameter
//
// From github.com/go-kit/kit/examples/profilesvc/transport.go
// ```
// ErrBadRouting is returned when an expected path variable is missing.
// It always indicates programmer error.`
// ```
var ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")


// Mounts service endpoints into an http.Handler
func MakeHTTPHandler(ctx context.Context, svc Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(svc)

	// These options are same for addsvc.. hmm
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST    /       posts by user ID a score, which is saved when greater than existing
	// GET     /:id    returns a score by id
	r.Methods("POST").Path("/").Handler(httptransport.NewServer(
		ctx,
		e.PostScoreEndpoint,
		decodePostScoreRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/{id}").Handler(httptransport.NewServer(
		ctx,
		e.PostScoreEndpoint,
		decodeGetScoreRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodePostScoreRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postScoreRequest // defined in endpoints
	if e := json.NewDecoder(r.Body).Decode(&req.Score); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetScoreRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getScoreRequest{Id: id}, nil
}

// implemented by all concrete response types that may contain errors
// used by response encoder
// TODO: read comment in endpoints.go
type error interface {
	error() error
}

// encodeResponse is the common method to encode all response types
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// business-logic error
		// rather than Go kit transport error
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// SKIP
// func encodeRequest(_ context.Context, req *http.Request, request interface{}) error

// encode json error response
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCodeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func statusCodeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
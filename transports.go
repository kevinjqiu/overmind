package overmind

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// ContentType defines the default content type for our responses
const ContentType = "application/json; charset=utf-8"

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil")
	}

	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(http.StatusBadRequest) // TODO: write different errors
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", ContentType)
	return json.NewEncoder(w).Encode(response)
}

func noopDecodeRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	return r, nil
}

func decodeGetZerglingByIDRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.New("Missing \"id\"")
	}
	request = getZerglingByIDRequest{id}
	return request, nil
}

func decodePostZerglingActionRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postZerglingActionRequest
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.New("Missing \"id\"")
	}

	if e := json.NewDecoder(r.Body).Decode(&req.command); e != nil {
		return nil, e
	}

	req.id = id
	return req, nil
}

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	endpoints := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("GET").Path("/_health").Handler(httptransport.NewServer(
		endpoints.GetHealthEndpoint,
		noopDecodeRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/zerglings/").Handler(httptransport.NewServer(
		endpoints.GetZerglingsEndpoint,
		noopDecodeRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/zerglings/").Handler(httptransport.NewServer(
		endpoints.PostZerglingsEndpoint,
		noopDecodeRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/zerglings/{id}").Handler(httptransport.NewServer(
		endpoints.GetZerglingByIDEndpoint,
		decodeGetZerglingByIDRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/zerglings/{id}").Handler(httptransport.NewServer(
		endpoints.PostZerglingActionEndpoint,
		decodePostZerglingActionRequest,
		encodeResponse,
		options...,
	))
	return r
}

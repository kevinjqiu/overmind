package overmind

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints collects all of the endpoints that compose an overmind service.
type Endpoints struct {
	GetHealthEndpoint       endpoint.Endpoint
	GetZerglingsEndpoint    endpoint.Endpoint
	PostZerglingsEndpoint   endpoint.Endpoint
	GetZerglingByIDEndpoint endpoint.Endpoint
}

type getHealthResponse struct {
	Health Health `json:"health,omitempty"`
	Err    error  `json:"err,omitempty"`
}

type getZerglingsResponse struct {
	Zerglings []Zergling `json:"zerglings,omitempty"`
	Err       error      `json:"err,omitempty"`
}

type postZerglingsResponse struct {
	Zergling Zergling `json:"zergling,omitempty"`
	Err      error    `json:"err,omitempty"`
}

type getZerglingByIDRequest struct {
	id string
}

type getZerglingByIDResponse struct {
	Zergling Zergling `json:"zergling,omitempty"`
	Err      error    `json:"err,omitempty"`
}

func makeGetHealthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		health, e := s.GetHealth(ctx)
		return getHealthResponse{
			Health: health,
			Err:    e,
		}, nil
	}
}

func makeGetZerglingsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		zerglings, e := s.GetZerglings(ctx)
		return getZerglingsResponse{
			Zerglings: zerglings,
			Err:       e,
		}, nil
	}
}

func makePostZerglingsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		zergling, e := s.PostZerglings(ctx)
		return postZerglingsResponse{Zergling: zergling, Err: e}, nil
	}
}

func makeGetZerglingByIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getZerglingByIDRequest)
		zergling, e := s.GetZerglingByID(ctx, req.id)
		return getZerglingByIDResponse{Zergling: zergling, Err: e}, nil
	}
}

// MakeServerEndpoints returns an Endpoint struct
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetHealthEndpoint:       makeGetHealthEndpoint(s),
		GetZerglingsEndpoint:    makeGetZerglingsEndpoint(s),
		PostZerglingsEndpoint:   makePostZerglingsEndpoint(s),
		GetZerglingByIDEndpoint: makeGetZerglingByIDEndpoint(s),
	}
}
